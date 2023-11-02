package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/files"
	"github.com/mjburtenshaw/macglab/glab"
	"github.com/mjburtenshaw/macglab/mrs"
	"github.com/mjburtenshaw/macglab/utils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	ApprovedFlag     bool
	BrowserFlag      bool
	DraftFlag        bool
	GroupFlag        bool
	ProjectsFlag     bool
	FlagAccessToken  string
	FlagGroupId      string
	FlagMe           int
	FlagUsernamesRaw string
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolVarP(&ApprovedFlag, "approved", "a", false, "Filter output to include MRs approved by the configured me user ID.")
	listCmd.PersistentFlags().BoolVarP(&BrowserFlag, "browser", "b", false, "Open merge requests in the browser.")
	listCmd.PersistentFlags().BoolVarP(&DraftFlag, "draft", "d", false, "Filter output to include draft merge requests.")
	listCmd.PersistentFlags().BoolVarP(&GroupFlag, "group", "g", false, "Filter output to the usernames configuration.")
	listCmd.PersistentFlags().StringVarP(&FlagGroupId, "group-id", "i", "", "Override the configured groud ID.")
	listCmd.PersistentFlags().IntVarP(&FlagMe, "me", "m", 0, "Override the configured me user ID.")
	listCmd.PersistentFlags().BoolVarP(&ProjectsFlag, "projects", "p", false, "Filter output to the projects configuration.")
	listCmd.PersistentFlags().StringVarP(&FlagAccessToken, "access-token", "t", "", "Override the configured access token.")
	listCmd.PersistentFlags().StringVarP(&FlagUsernamesRaw, "users", "u", "", "Filter output to the specified usernames.")
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
	- Use the '-t, --access-token' flag to override the configured access token.
	- Use the '-u, --users' flag to override configured usernames and only filter on usernames you provided. Accepts a CSV string of usernames.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Read(files.MacglabConfigUrl)
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		accessToken, groupId, me, flagUsernames, shouldAskToUpdateAccessToken, shouldAskToUpdateGroupId, shouldAskToUpdateMe := parseFlags(conf)

		glabClient, err := glab.Initialize(accessToken)
		if err != nil {
			log.Fatalf("Failed to initialize gitlab client: %v", err)
		}

		allMrs, err := fetchMergeRequests(glabClient, conf, groupId, me, &DraftFlag, &GroupFlag, &ProjectsFlag, flagUsernames)
		if err != nil {
			log.Fatalf("Failed to fetch merge requests: %v", err)
		}

		mrs.PrintMergeRequests(allMrs)

		if BrowserFlag {
			if err := mrs.OpenMergeRequests(allMrs); err != nil {
				log.Printf("Failed to open merge requests in the browser: %v", err)
			}
		}

		if shouldAskToUpdateAccessToken {
			response := utils.AskBinaryQuestion("Do you want to use the same access token in the future? (yes/no): ")
			if strings.HasPrefix(strings.ToLower(response), "y") {
				config.Update(files.MacglabConfigUrl, "access_token", FlagAccessToken)
			}
		}

		if shouldAskToUpdateGroupId {
			response := utils.AskBinaryQuestion("Do you want to use the same group ID in the future? (yes/no): ")
			if strings.HasPrefix(strings.ToLower(response), "y") {
				config.Update(files.MacglabConfigUrl, "group_id", FlagGroupId)
			}
		}

		if shouldAskToUpdateMe {
			response := utils.AskBinaryQuestion("Do you want to use the same me user ID in the future? (yes/no): ")
			if strings.HasPrefix(strings.ToLower(response), "y") {
				config.Update(files.MacglabConfigUrl, "me", fmt.Sprintf("%d", FlagMe))
			}
		}
	},
}

func parseFlags(conf *config.Config) (accessToken string, groupId string, me int, flagUsernames []string, shouldAskToUpdateAccessToken bool, shouldAskToUpdateGroupId bool, shouldAskToUpdateMe bool) {
	accessToken = conf.AccessToken
	shouldAskToUpdateAccessToken = false
	if FlagAccessToken != "" {
		accessToken = FlagAccessToken
		shouldAskToUpdateAccessToken = true
	}

	groupId = conf.GroupId
	shouldAskToUpdateGroupId = false
	if FlagGroupId != "" {
		groupId = FlagGroupId
		shouldAskToUpdateGroupId = true
	}

	me = conf.Me
	shouldAskToUpdateMe = false
	if FlagMe != 0 {
		me = FlagMe
		shouldAskToUpdateMe = true
	}

	FlagUsernamesRaw = strings.ReplaceAll(FlagUsernamesRaw, " ", "")
	if FlagUsernamesRaw != "" {
		flagUsernames = strings.Split(FlagUsernamesRaw, ",")
	}

	return accessToken, groupId, me, flagUsernames, shouldAskToUpdateAccessToken, shouldAskToUpdateGroupId, shouldAskToUpdateMe
}

func fetchMergeRequests(glabClient *glab.TGitlabClient, conf *config.Config, groupId string, me int, DraftFlag, GroupFlag, ProjectsFlag *bool, flagUsernames []string) ([]*gitlab.MergeRequest, error) {
	var allMrs []*gitlab.MergeRequest

	if (!*GroupFlag && !*ProjectsFlag) || *GroupFlag {
		usernames := chooseUsernames(flagUsernames, conf.Usernames)
		groupMrs, err := mrs.FetchGroupMergeRequests(glabClient, groupId, usernames, DraftFlag)
		if err != nil {
			return nil, err
		}
		allMrs = append(allMrs, groupMrs...)
	}

	if (!*GroupFlag && !*ProjectsFlag) || *ProjectsFlag {
		allProjectUsernames := conf.Projects["all"]

		for project, thisProjectUsernames := range conf.Projects {
			if project != "all" {
				projectUsernames := append(thisProjectUsernames, allProjectUsernames...)
				usernames := chooseUsernames(flagUsernames, projectUsernames)
				projectMrs, err := mrs.FetchProjectMergeRequests(glabClient, project, usernames, DraftFlag)
				if err != nil {
					return nil, err
				}
				allMrs = append(allMrs, projectMrs...)
			}
		}
	}

	allMrs = dedupeMergeRequests(allMrs)

	if !ApprovedFlag && me != 0 {
		mrsNotApprovedByMe, err := excludeMrsApprovedByMe(glabClient, groupId, me, allMrs)
		if err != nil {
			return nil, err
		}
		allMrs = mrsNotApprovedByMe
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
	approvedMrs, err := mrs.GetMergeRequestsApprovedByMe(glabClient, groupId, me, &DraftFlag)
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
