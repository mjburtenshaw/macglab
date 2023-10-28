package cmd

import (
	"log"
	"strings"

	"github.com/mjburtenshaw/macglab/config"
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
	FlagAccessToken	 string
	FlagUsernamesRaw string
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().BoolVarP(&ApprovedFlag, "approved", "a", false, "Filter output to include MRs approved by the configured ME user ID.")
	listCmd.PersistentFlags().BoolVarP(&BrowserFlag, "browser", "b", false, "Open merge requests in the browser.")
	listCmd.PersistentFlags().BoolVarP(&DraftFlag, "draft", "d", false, "Filter output to include draft merge requests.")
	listCmd.PersistentFlags().BoolVarP(&GroupFlag, "group", "g", false, "Filter output to the usernames configuration.")
	listCmd.PersistentFlags().BoolVarP(&ProjectsFlag, "projects", "p", false, "Filter output to the projects configuration.")
	listCmd.PersistentFlags().StringVarP(&FlagAccessToken, "access-token", "t", "", "Override the configured access token.")
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
	- Use the '-t, --access-token' flag to override the configured access token.
	- Use the '-u, --users' flag to override configured usernames and only filter on usernames you provided. Accepts a CSV string of usernames.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Read(config.MacglabConfigUrl)
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		accessToken := conf.AccessToken
		shouldAskToUpdateAccessToken := false
		if (FlagAccessToken != "") {
			accessToken = FlagAccessToken
			shouldAskToUpdateAccessToken = true
		}

		err = glab.Initialize(accessToken)
		if err != nil {
			log.Fatalf("Failed to initialize gitlab client: %v", err)
		}

		FlagUsernamesRaw = strings.ReplaceAll(FlagUsernamesRaw, " ", "")
		var flagUsernames []string
		if FlagUsernamesRaw != "" {
			flagUsernames = strings.Split(FlagUsernamesRaw, ",")
		}

		allMrs, err := fetchMergeRequests(conf, &DraftFlag, &GroupFlag, &ProjectsFlag, flagUsernames)
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
				config.Update(config.MacglabConfigUrl, "ACCESS_TOKEN", FlagAccessToken)
			}
		}
	},
}

func fetchMergeRequests(conf *config.Config, DraftFlag, GroupFlag, ProjectsFlag *bool, flagUsernames []string) ([]*gitlab.MergeRequest, error) {
	var allMrs []*gitlab.MergeRequest

	if (!*GroupFlag && !*ProjectsFlag) || *GroupFlag {
		usernames := chooseUsernames(flagUsernames, conf.Usernames)
		groupMrs, err := mrs.FetchGroupMergeRequests(conf.GroupId, usernames, DraftFlag)
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
				projectMrs, err := mrs.FetchProjectMergeRequests(project, usernames, DraftFlag)
				if err != nil {
					return nil, err
				}
				allMrs = append(allMrs, projectMrs...)
			}
		}
	}

	allMrs = dedupeMergeRequests(allMrs)

	if !ApprovedFlag && conf.Me != 0 {
		mrsNotApprovedByMe, err := excludeMrsApprovedByMe(conf, allMrs)
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

func excludeMrsApprovedByMe(conf *config.Config, allMrs []*gitlab.MergeRequest) ([]*gitlab.MergeRequest, error) {
	approvedMrs, err := mrs.GetMergeRequestsApprovedByMe(conf.GroupId, conf.Me, &DraftFlag)
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
