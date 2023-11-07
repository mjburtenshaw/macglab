package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/files"
	"github.com/mjburtenshaw/macglab/glab"
	"github.com/mjburtenshaw/macglab/mrs"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ResolvedFlags struct {
    accessToken   string
    groupId       string
    me            int
    flagUsernames []string
}

type TrueUpFlags map[string]bool

var (
    approvedFlag     bool
    browserFlag      bool
    draftFlag        bool
    groupFlag        bool
    projectsFlag     bool
    readyFlag        bool
    flagAccessToken  string
    flagGroupId      string
    flagMe           int
    flagUsernamesRaw string
)

func init() {
    rootCmd.AddCommand(listCmd)
    listCmd.PersistentFlags().BoolVarP(&approvedFlag, "approved", "a", false, "Filter output to include MRs approved by the configured me user ID.")
    listCmd.PersistentFlags().BoolVarP(&browserFlag, "browser", "b", false, "Open merge requests in the browser.")
    listCmd.PersistentFlags().BoolVarP(&draftFlag, "draft", "d", false, "Filter output to include draft merge requests.")
    listCmd.PersistentFlags().BoolVarP(&groupFlag, "group", "g", false, "Filter output to the usernames configuration.")
    listCmd.PersistentFlags().StringVarP(&flagGroupId, "group-id", "i", "", "Override the configured groud ID.")
    listCmd.PersistentFlags().IntVarP(&flagMe, "me", "m", 0, "Override the configured me user ID.")
    listCmd.PersistentFlags().BoolVarP(&projectsFlag, "projects", "p", false, "Filter output to the projects configuration.")
    listCmd.PersistentFlags().BoolVarP(&readyFlag, "ready", "r", false, "Filter output to include merge requests that are ready to merge.")
    listCmd.PersistentFlags().StringVarP(&flagAccessToken, "access-token", "t", "", "Override the configured access token.")
    listCmd.PersistentFlags().StringVarP(&flagUsernamesRaw, "users", "u", "", "Filter output to the specified usernames.")
}

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List merge requests",
    Long: `List merge requests using the following options:
    - Use the '-a, --approved' flag to filter output to include MRs approved by the configured 'me' user ID.
    - Use the '-b, --browser' flag to open MRs in the browser.
    - Use the '-d, --drafts' flag to include draft MRs.
    - Use the '-g, --group' flag to filter output to the usernames configuration.
    - Use the '-i, --group-id' flag to override the configured group ID.
    - Use the '-m, --me' flag to override the configured me user ID.
    - Use the '-p, --projects' flag to filter output to the projects configuration.
    - Use the '-r, --ready' flag to filter output to include merge requests that are ready to merge.
    - Use the '-t, --access-token' flag to override the configured access token.
    - Use the '-u, --users' flag to override configured usernames and only filter on usernames you provided. Accepts a CSV string of usernames.`,
    Run: func(cmd *cobra.Command, args []string) {
        conf, err := config.Read(files.MacglabConfigUrl)
        if err != nil {
            log.Fatalf("Failed to read config: %v", err)
        }

        resolvedFlags, trueUpFlags := resolveFlags(conf)

        glabClient, err := glab.Initialize(resolvedFlags.accessToken)
        if err != nil {
            log.Fatalf("Failed to initialize gitlab client: %v", err)
        }

        allMrs, err := fetchMergeRequests(glabClient, conf, resolvedFlags.groupId, resolvedFlags.me, &draftFlag, &groupFlag, &projectsFlag, resolvedFlags.flagUsernames)
        if err != nil {
            log.Fatalf("Failed to fetch merge requests: %v", err)
        }

        mrs.PrintMergeRequests(allMrs)

        if browserFlag {
            if err := mrs.OpenMergeRequests(allMrs); err != nil {
                log.Printf("Failed to open merge requests in the browser: %v", err)
            }
        }

        config.TrueUp([]config.TrueUpKit{
            {
                ShouldAsk:  trueUpFlags["shouldAskToUpdateAccessToken"],
                Question:   "Do you want to use the same access token in the future? (yes/no): ",
                ConfigAttr: "access_token",
                NextValue:  flagAccessToken,
            },
            {
                ShouldAsk:  trueUpFlags["shouldAskToUpdateGroupId"],
                Question:   "Do you want to use the same group ID in the future? (yes/no): ",
                ConfigAttr: "group_id",
                NextValue:  flagGroupId,
            },
            {
                ShouldAsk:  trueUpFlags["shouldAskToUpdateMe"],
                Question:   "Do you want to use the same me user ID in the future? (yes/no): ",
                ConfigAttr: "me",
                NextValue:  fmt.Sprintf("%d", flagMe),
            },
        })
    },
}

func resolveFlags(conf *config.Config) (resolvedFlags ResolvedFlags, trueUpFlags TrueUpFlags) {
    resolvedFlags = ResolvedFlags{
        accessToken:   conf.AccessToken,
        groupId:       conf.GroupId,
        me:            conf.Me,
        flagUsernames: []string{},
    }

    trueUpFlags = make(TrueUpFlags)

    if flagAccessToken != "" {
        resolvedFlags.accessToken = flagAccessToken
        trueUpFlags["shouldAskToUpdateAccessToken"] = true
    }

    if flagGroupId != "" {
        resolvedFlags.groupId = flagGroupId
        trueUpFlags["shouldAskToUpdateGroupId"] = true
    }

    if flagMe != 0 {
        resolvedFlags.me = flagMe
        trueUpFlags["shouldAskToUpdateMe"] = true
    }

    flagUsernamesRaw = strings.ReplaceAll(flagUsernamesRaw, " ", "")
    if flagUsernamesRaw != "" {
        resolvedFlags.flagUsernames = strings.Split(flagUsernamesRaw, ",")
    }

    return resolvedFlags, trueUpFlags
}

func fetchMergeRequests(glabClient *glab.TGitlabClient, conf *config.Config, groupId string, me int, draftFlag, groupFlag, projectsFlag *bool, flagUsernames []string) ([]*gitlab.MergeRequest, error) {
    var allMrs []*gitlab.MergeRequest

    if (!*groupFlag && !*projectsFlag) || *groupFlag {
        usernames := chooseUsernames(flagUsernames, conf.Usernames)
        groupMrs, err := mrs.FetchGroupMergeRequests(glabClient, groupId, usernames, draftFlag)
        if err != nil {
            return nil, err
        }
        allMrs = append(allMrs, groupMrs...)
    }

    if (!*groupFlag && !*projectsFlag) || *projectsFlag {
        allProjectUsernames := conf.Projects["all"]

        for project, thisProjectUsernames := range conf.Projects {
            if project != "all" {
                projectUsernames := append(thisProjectUsernames, allProjectUsernames...)
                usernames := chooseUsernames(flagUsernames, projectUsernames)
                projectMrs, err := mrs.FetchProjectMergeRequests(glabClient, project, usernames, draftFlag)
                if err != nil {
                    return nil, err
                }
                allMrs = append(allMrs, projectMrs...)
            }
        }
    }

    mrsInReviewByMe, err := mrs.FetchReviewerMergeRequests(glabClient, groupId, me, draftFlag)
    if err != nil {
        return nil, err
    }
    allMrs = append(allMrs, mrsInReviewByMe...)

    allMrs = dedupeMergeRequests(allMrs)

    if !approvedFlag && me != 0 {
        mrsNotApprovedByMe, err := excludeMrsApprovedByMe(glabClient, groupId, me, allMrs)
        if err != nil {
            return nil, err
        }
        allMrs = mrsNotApprovedByMe
    }

    // Filter out MRs that are ready to merge.
    // See https://docs.gitlab.com/ee/api/merge_requests.html#merge-status for a list of statuses.
    // [The `go-gitlab` maintainer believes adding enum support requires too much maintenance](https://github.com/xanzy/go-gitlab/pull/1774#issuecomment-1728723321).
    if !readyFlag {
        mrsNotReadyToMerge := []*gitlab.MergeRequest{}
        for _, mr := range allMrs {
            if mr.DetailedMergeStatus != "mergeable" {
                mrsNotReadyToMerge = append(mrsNotReadyToMerge, mr)
            }
        }
        allMrs = mrsNotReadyToMerge
    }

    return allMrs, nil
}

// chooseUsernames chooses usernames provided via the user flag over the config.
func chooseUsernames(flagUsernames []string, configUsernames []string) []string {
    if len(flagUsernames) != 0 {
        return flagUsernames
    }
    return configUsernames
}

func dedupeMergeRequests(mergeRequests []*gitlab.MergeRequest) []*gitlab.MergeRequest {
    seen := map[string]bool{}
    result := []*gitlab.MergeRequest{}

    for _, mergeRequest := range mergeRequests {

        if !seen[mergeRequest.WebURL] {
            seen[mergeRequest.WebURL] = true
            result = append(result, mergeRequest)
        }
    }

    return result
}

func excludeMrsApprovedByMe(glabClient *glab.TGitlabClient, groupId string, me int, allMrs []*gitlab.MergeRequest) ([]*gitlab.MergeRequest, error) {
    approvedMrs, err := mrs.GetMergeRequestsApprovedByMe(glabClient, groupId, me, &draftFlag)
    if err != nil {
        return nil, err
    }

    mrsNotApprovedByMe := []*gitlab.MergeRequest{}
    for _, mr := range allMrs {
        isApproved := false
        for _, approvedMr := range approvedMrs {
            if mr.IID == approvedMr.IID {
                isApproved = true
                break
            }
        }
        if !isApproved {
            mrsNotApprovedByMe = append(mrsNotApprovedByMe, mr)
        }
    }

    return mrsNotApprovedByMe, nil
}
