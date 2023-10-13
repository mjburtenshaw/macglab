package main

import (
	"flag"
	"log"
	"strings"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/mrs"
	"github.com/xanzy/go-gitlab"
)

func main() {
	browserFlag := flag.Bool("browser", false, "Open merge requests in the browser.")
	draftFlag := flag.Bool("draft", false, "Filter output to include draft merge requests.")
	groupFlag := flag.Bool("group", false, "Filter output to the usernames configuration.")
	projectsFlag := flag.Bool("projects", false, "Filter output to the projects configuration.")

	var flagUsernamesRaw string
	flag.StringVar(&flagUsernamesRaw, "users", "", "Filter output to the specified usernames.")

	flag.Parse()

	flagUsernamesRaw = strings.ReplaceAll(flagUsernamesRaw, " ", "")

	flagUsernames := strings.Split(flagUsernamesRaw, ",")

	allMrs, err := fetchMergeRequests(draftFlag, groupFlag, projectsFlag, flagUsernames)
	if err != nil {
		log.Fatalf("Failed to fetch merge requests: %v", err)
	}

	mrs.PrintMergeRequests(allMrs)

	if *browserFlag {
		if err := mrs.OpenMergeRequests(allMrs); err != nil {
			log.Printf("Failed to open merge requests in browser: %v", err)
		}
	}
}

func fetchMergeRequests(draftFlag, groupFlag, projectsFlag *bool, flagUsernames []string) ([]*gitlab.MergeRequest, error) {
	var allMrs []*gitlab.MergeRequest

	if (!*groupFlag && !*projectsFlag) || *groupFlag {
		usernames := chooseUsernames(flagUsernames, config.Usernames)
		groupMrs, err := mrs.FetchGroupMergeRequests(usernames, draftFlag)
		if err != nil {
			return nil, err
		}
		allMrs = append(allMrs, groupMrs...)
	}

	if (!*groupFlag && !*projectsFlag) || *projectsFlag {
		allProjectUsernames := config.Projects["all"]

		for project, thisProjectUsernames := range config.Projects {
			if project != "all" {
				projectUsernames := append(thisProjectUsernames, allProjectUsernames...)
				usernames := chooseUsernames(flagUsernames, projectUsernames)
				projectMrs, err := mrs.FetchProjectMergeRequests(project, usernames, draftFlag)
				if err != nil {
					return nil, err
				}
				allMrs = append(allMrs, projectMrs...)
			}
		}
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
