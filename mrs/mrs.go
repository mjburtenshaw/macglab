package mrs

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/xanzy/go-gitlab"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/glab"
)

func openURL(url string) error {
	var cmd string
	var args []string
	
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	
	return exec.Command(cmd, args...).Start()
}

// getWip converts a boolean pointer to a string representation used for querying GitLab's
// Merge Request API with regard to draft (Work In Progress - WIP) status.
// GitLab currently does not support a `draft` query parameter, instead, it uses `wip` (Work In Progress).
//
// The function behaves as follows:
// - If `shouldIncludeDrafts` is nil or points to `false`, "no" is returned to exclude drafts from results.
// - If `shouldIncludeDrafts` points to `true`, an empty string is returned, which does not filter results
//   based on draft status, effectively including drafts in the results.
//
// Example:
//   drafts := true
//   queryParam := getWip(&drafts)  // queryParam will be ""
//
// GitLab Merge Requests API documentation:
// https://docs.gitlab.com/ee/api/merge_requests.html#list-merge-requests
//
// Usage:
// The returned string is intended to be used as the value for the `wip` query parameter in
// GitLab's Merge Requests API.
//
//   options := &gitlab.ListMergeRequestsOptions{
//       WIP: gitlab.String(getWip(&drafts)),
//   }
//   mergeRequests, _, err := gitlabClient.MergeRequests.ListMergeRequests(projectID, options)
func getWip(shouldIncludeDrafts *bool) string {
	wip := "no"

	if shouldIncludeDrafts != nil && *shouldIncludeDrafts {
		wip = ""
	}

	return wip
}

func FetchGroupMergeRequests(shouldIncludeDrafts *bool) [] *gitlab.MergeRequest {
	var groupMrs [] *gitlab.MergeRequest

	for _, username := range config.Usernames {
		userMrs, _, err := glab.Client.MergeRequests.ListGroupMergeRequests(config.GroupId, &gitlab.ListGroupMergeRequestsOptions{
			AuthorUsername: gitlab.String(username),
			State: gitlab.String("opened"),
			WIP: gitlab.String(getWip(shouldIncludeDrafts)),
		})
		if err != nil {
			log.Fatalf("💀 Failed to get merge request for %s: %v", username, err)
		}

		groupMrs = append(groupMrs, userMrs...)
	}
	
	return groupMrs
}

func FetchProjectMergeRequests(projectId string, usernames []string, shouldIncludeDrafts *bool) [] *gitlab.MergeRequest {
	var projectMrs [] *gitlab.MergeRequest
	
	for _, username := range usernames {
		userMrs, _, err := glab.Client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			AuthorUsername: gitlab.String(username),
			State: gitlab.String("opened"),
			WIP: gitlab.String(getWip(shouldIncludeDrafts)),
		})
		if err != nil {
			log.Fatalf("💀 Failed to get merge request for %s: %v", username, err)
		}

		projectMrs = append(projectMrs, userMrs...)
	}
	
	return projectMrs
}

func PrintMergeRequests(mrs []*gitlab.MergeRequest) {
	for _, mr := range mrs {
		fmt.Printf("@%s: %s\n", mr.Author.Username, mr.WebURL)
	}
}

func OpenMergeRequests(mrs []*gitlab.MergeRequest) {
	for _, mr := range mrs {
		openURL(mr.WebURL)
	}
}
