package main

import (
	"flag"

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

	var allMrs []*gitlab.MergeRequest

	if (!*groupFlag && !*projectsFlag) || *groupFlag {
		groupMrs := mrs.FetchGroupMergeRequests(draftFlag)
		allMrs = append(allMrs, groupMrs...)
	}

	if (!*groupFlag && !*projectsFlag) || *projectsFlag {
		allUsernames := config.Projects["all"]

		for project, usernames := range config.Projects {
				if project != "all" {
						combinedUsernames := append(usernames, allUsernames...)
						projectMrs := mrs.FetchProjectMergeRequests(project, combinedUsernames, draftFlag)
						allMrs = append(allMrs, projectMrs...)
				}
		}
	}

	mrs.PrintMergeRequests(allMrs)

	if *browserFlag {
		mrs.OpenMergeRequests(allMrs)
	}
}
