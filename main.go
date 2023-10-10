package main

import (
	"flag"
	"log"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/mrs"
	"github.com/xanzy/go-gitlab"
)

func main() {
	browserFlag := flag.Bool("browser", false, "Open merge requests in the browser.")
	draftFlag := flag.Bool("draft", false, "Filter output to include draft merge requests.")
	groupFlag := flag.Bool("group", false, "Filter output to the usernames configuration.")
	projectsFlag := flag.Bool("projects", false, "Filter output to the projects configuration.")

	flag.Parse()

	allMrs, err := fetchMergeRequests(draftFlag, groupFlag, projectsFlag)
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

func fetchMergeRequests(draftFlag, groupFlag, projectsFlag *bool) ([]*gitlab.MergeRequest, error) {
	var allMrs []*gitlab.MergeRequest

	if (!*groupFlag && !*projectsFlag) || *groupFlag {
		groupMrs, err := mrs.FetchGroupMergeRequests(draftFlag)
		if err != nil {
			return nil, err
		}
		allMrs = append(allMrs, groupMrs...)
	}

	if (!*groupFlag && !*projectsFlag) || *projectsFlag {
		allUsernames := config.Projects["all"]

		for project, usernames := range config.Projects {
			if project != "all" {
				combinedUsernames := append(usernames, allUsernames...)
				projectMrs, err := mrs.FetchProjectMergeRequests(project, combinedUsernames, draftFlag)
				if err != nil {
					return nil, err
				}
				allMrs = append(allMrs, projectMrs...)
			}
		}
	}

	return allMrs, nil
}
