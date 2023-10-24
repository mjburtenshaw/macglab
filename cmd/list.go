package cmd

import (
	"log"
	"strings"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/glab"
	"github.com/mjburtenshaw/macglab/mrs"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	ApprovedFlag     bool
	BrowserFlag      bool
	DraftFlag        bool
	GroupFlag        bool
	ProjectsFlag     bool
	FlagUsernamesRaw string
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolVarP(&ApprovedFlag, "approved", "a", false, "Filter output to include MRs approved by the configured ME user ID.")
	listCmd.PersistentFlags().BoolVarP(&BrowserFlag, "browser", "b", false, "Open merge requests in the browser.")
	listCmd.PersistentFlags().BoolVarP(&DraftFlag, "draft", "d", false, "Filter output to include draft merge requests.")
	listCmd.PersistentFlags().BoolVarP(&GroupFlag, "group", "g", false, "Filter output to the usernames configuration.")
	listCmd.PersistentFlags().BoolVarP(&ProjectsFlag, "projects", "p", false, "Filter output to the projects configuration.")
	listCmd.PersistentFlags().StringVarP(&FlagUsernamesRaw, "users", "u", "", "Filter output to the specified usernames.")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List merge requests",
	Long: `List merge requests using the following options:
	- Use the '-a, --approved' flag to filter output to include MRs approved by the configured 'ME' user ID.
	- Use the '-b, --browser' flag to open MRs in the browser.
	- Use the '-d, --drafts' flag to include draft MRs.
	- Use the '-g, --group' flag to filter output to the usernames configuration.
	- Use the '-p, --projects' flag to filter output to the projects configuration.
	- Use the '-u, --users' flag to override configured usernames and only filter on usernames you provided. Accepts a CSV string of usernames.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Read()
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		err = glab.Initialize()
		if err != nil {
			log.Fatalf("Failed to initialize gitlab client: %v", err)
		}

		FlagUsernamesRaw = strings.ReplaceAll(FlagUsernamesRaw, " ", "")
		var flagUsernames []string
		if FlagUsernamesRaw != "" {
			flagUsernames = strings.Split(FlagUsernamesRaw, ",")
		}

		allMrs, err := fetchMergeRequests(&DraftFlag, &GroupFlag, &ProjectsFlag, flagUsernames)
		if err != nil {
			log.Fatalf("Failed to fetch merge requests: %v", err)
		}

		mrs.PrintMergeRequests(allMrs)

		if BrowserFlag {
			if err := mrs.OpenMergeRequests(allMrs); err != nil {
				log.Printf("Failed to open merge requests in the browser: %v", err)
			}
		}
	},
}

func fetchMergeRequests(DraftFlag, GroupFlag, ProjectsFlag *bool, flagUsernames []string) ([]*gitlab.MergeRequest, error) {
	var allMrs []*gitlab.MergeRequest

	if (!*GroupFlag && !*ProjectsFlag) || *GroupFlag {
		usernames := chooseUsernames(flagUsernames, config.Usernames)
		groupMrs, err := mrs.FetchGroupMergeRequests(usernames, DraftFlag)
		if err != nil {
			return nil, err
		}
		allMrs = append(allMrs, groupMrs...)
	}

	if (!*GroupFlag && !*ProjectsFlag) || *ProjectsFlag {
		allProjectUsernames := config.Projects["all"]

		for project, thisProjectUsernames := range config.Projects {
			if project != "all" {
				projectUsernames := append(thisProjectUsernames, allProjectUsernames...)
				usernames := chooseUsernames(flagUsernames, projectUsernames)
				projectMrs, err := mrs.FetchProjectMergeRequests(project, usernames, DraftFlag)
				if err != nil {
					return nil, err
				}
				allMrs = append(allMrs, projectMrs...)
			}
		}
	}

	allMrs = dedupeMergeRequests(allMrs)

	if !ApprovedFlag && config.Me != 0 {
		mrsNotApprovedByMe, err := excludeMrsApprovedByMe(allMrs)
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

func excludeMrsApprovedByMe(allMrs []*gitlab.MergeRequest) ([]*gitlab.MergeRequest, error) {
	approvedMrs, err := mrs.GetMergeRequestsApprovedByMe(config.Me, &DraftFlag)
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
