package cmd

import (
	"fmt"
	"log"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/files"
	"github.com/mjburtenshaw/macglab/flags"
	"github.com/mjburtenshaw/macglab/glab"
	"github.com/mjburtenshaw/macglab/mrs"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func init() {
	rootCmd.AddCommand(listCmd)
	flags.AddListFlags(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List merge requests",
	Long: `list

Prints GitLab Merge Request (MRs) authors and URLs to the terminal.

list fetches MRs meeting ALL the following criteria:
- State is open.
- Belongs to the configured group ID.
- Is NOT a draft.
- Meets ANY of the following criteria:
	- The author is listed in the configured usernames.
	- The author is listed in ANY of the configured projects; but it only returns MRs for projects the author is listed under.
	- You are listed as a reviewer.

list then excludes MRs meeting the following criteria:
- Approved by you.
- Mergeable MRs where you are NOT the author.

Note: group and projects are not mutually exclusive. If neither are provided, the program will run as if both are provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Read(files.MacglabConfigUrl)
		if err != nil {
			log.Printf("Failed to read config: %v", err)
            return
		}

		listFlags := flags.GetListFlags(conf)

		glabClient, err := glab.Initialize(listFlags.Resolved.AccessToken)
		if err != nil {
			log.Printf("Failed to initialize gitlab client: %v", err)
            return
		}

		allMrs, err := fetchMergeRequests(glabClient, conf, listFlags.Resolved, listFlags.Boolean)
		if err != nil {
			log.Printf("Failed to fetch merge requests: %v", err)
            return
		}

		if listFlags.Boolean.Count {
			fmt.Printf("count: %v\n", len(allMrs))
		}

		mrs.PrintMergeRequests(allMrs)

		if listFlags.Boolean.Browser {
			if err := mrs.OpenMergeRequests(allMrs); err != nil {
				log.Printf("Failed to open merge requests in the browser: %v", err)
			}
		}

		config.TrueUp([]config.TrueUpKit{
			{
				ShouldAsk:  listFlags.TrueUp["shouldAskToUpdateAccessToken"],
				Question:   "Do you want to use the same access token in the future? (yes/no): ",
				ConfigAttr: "access_token",
				NextValue:  listFlags.RawValue.AccessToken,
			},
			{
				ShouldAsk:  listFlags.TrueUp["shouldAskToUpdateGroupId"],
				Question:   "Do you want to use the same group ID in the future? (yes/no): ",
				ConfigAttr: "group_id",
				NextValue:  listFlags.RawValue.GroupId,
			},
			{
				ShouldAsk:  listFlags.TrueUp["shouldAskToUpdateMe"],
				Question:   "Do you want to use the same me user ID in the future? (yes/no): ",
				ConfigAttr: "me",
				NextValue:  fmt.Sprintf("%d", listFlags.RawValue.Me),
			},
		})
	},
}

func fetchMergeRequests(glabClient *glab.TGitlabClient, conf *config.Config, resolvedFlags flags.ResolvedFlags, booleanFlags flags.BooleanFlags) ([]*gitlab.MergeRequest, error) {
	var allMrs []*gitlab.MergeRequest

	if (!booleanFlags.Group && !booleanFlags.Projects) || booleanFlags.Group {
		usernames := chooseUsernames(resolvedFlags.Usernames, conf.Usernames)
		groupMrs, err := mrs.FetchGroupMergeRequests(glabClient, resolvedFlags.GroupId, usernames, &booleanFlags.Draft)
		if err != nil {
			return nil, err
		}
		allMrs = append(allMrs, groupMrs...)
	}

	if (!booleanFlags.Group && !booleanFlags.Projects) || booleanFlags.Projects {
		allProjectUsernames := conf.Projects["all"]

		for project, thisProjectUsernames := range conf.Projects {
			if project != "all" {
				projectUsernames := append(thisProjectUsernames, allProjectUsernames...)
				usernames := chooseUsernames(resolvedFlags.Usernames, projectUsernames)
				projectMrs, err := mrs.FetchProjectMergeRequests(glabClient, project, usernames, &booleanFlags.Draft)
				if err != nil {
					return nil, err
				}
				allMrs = append(allMrs, projectMrs...)
			}
		}
	}

	mrsInReviewByMe, err := mrs.FetchReviewerMergeRequests(glabClient, resolvedFlags.GroupId, resolvedFlags.Me, &booleanFlags.Draft)
	if err != nil {
		return nil, err
	}
	allMrs = append(allMrs, mrsInReviewByMe...)

	allMrs = dedupeMergeRequests(allMrs)

	if !booleanFlags.Approved && resolvedFlags.Me != 0 {
		mrsNotApprovedByMe, err := excludeMrsApprovedByMe(glabClient, resolvedFlags.GroupId, resolvedFlags.Me, booleanFlags.Draft, allMrs)
		if err != nil {
			return nil, err
		}
		allMrs = mrsNotApprovedByMe
	}

	// Filter out MRs that are ready to merge, unless the given `me` GitLab user ID is the author.
	// See https://docs.gitlab.com/ee/api/merge_requests.html#merge-status for a list of statuses.
	// [The `go-gitlab` maintainer believes adding enum support requires too much maintenance](https://github.com/xanzy/go-gitlab/pull/1774#issuecomment-1728723321).
	if !booleanFlags.Ready {
		mrsNotReadyToMerge := []*gitlab.MergeRequest{}
		for _, mr := range allMrs {
			if mr.DetailedMergeStatus != "mergeable"  || mr.Author.ID == resolvedFlags.Me {
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

func excludeMrsApprovedByMe(glabClient *glab.TGitlabClient, groupId string, me int, draftFlag bool, allMrs []*gitlab.MergeRequest) ([]*gitlab.MergeRequest, error) {
	approvedMrs, err := mrs.GetMergeRequestsApprovedByMe(glabClient, groupId, me, &draftFlag)
	if err != nil {
		return nil, err
	}

	mrsNotApprovedByMe := []*gitlab.MergeRequest{}
	for _, mr := range allMrs {
		isApproved := false
		for _, approvedMr := range approvedMrs {
			if mr.ID == approvedMr.ID {
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
