package mrs

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/mjburtenshaw/macglab/glab"
	"github.com/xanzy/go-gitlab"
)

// openURL opens the specified URL in the user's default browser.
func openURL(url string) error {
	if url == "" {
		return errors.New("url cannot be empty")
	}

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

// getWIPQueryParamPointer converts a boolean pointer to a string representation for querying GitLab's
// Merge Request API with regard to Work In Progress (WIP) status.
func getWIPQueryParamPointer(shouldIncludeDrafts *bool) *string {
	// We don't want to filter by WIP status if we're including drafts.
	// Setting wip to "yes" will filter results by *only* drafts.
	if shouldIncludeDrafts == nil || *shouldIncludeDrafts {
		return nil
	}
	wip := "no"
	return &wip
}

// FetchGroupMergeRequests fetches merge requests for a group from GitLab.
func FetchGroupMergeRequests(glabClient *glab.TGitlabClient, groupId string, usernames []string, shouldIncludeDrafts *bool) ([]*gitlab.MergeRequest, error) {
	var groupMrs []*gitlab.MergeRequest

	for _, username := range usernames {
		userMrs, err := fetchUserMergeRequests(glabClient, groupId, username, shouldIncludeDrafts)
		if err != nil {
			return nil, err
		}
		groupMrs = append(groupMrs, userMrs...)
	}

	return groupMrs, nil
}

// fetchUserMergeRequests fetches merge requests for a specific user within a group from GitLab.
func fetchUserMergeRequests(glabClient *glab.TGitlabClient, groupId string, username string, shouldIncludeDrafts *bool) ([]*gitlab.MergeRequest, error) {
	userMrs, _, err := glabClient.MergeRequests.ListGroupMergeRequests(groupId, &gitlab.ListGroupMergeRequestsOptions{
		AuthorUsername: gitlab.String(username),
		State:          gitlab.String("opened"),
		WIP:            getWIPQueryParamPointer(shouldIncludeDrafts),
	})
	if err != nil {
		log.Printf("Failed to get merge requests for %s: %v\n", username, err)
		return nil, err
	}

	return userMrs, nil
}

// FetchUserMergeRequests fetches merge requests for a specific reviewer within a group from GitLab.
func FetchReviewerMergeRequests(glabClient *glab.TGitlabClient, groupId string, userId int, shouldIncludeDrafts *bool) ([]*gitlab.MergeRequest, error) {
	userMrs, _, err := glabClient.MergeRequests.ListGroupMergeRequests(groupId, &gitlab.ListGroupMergeRequestsOptions{
		ReviewerID: 		gitlab.ReviewerID(userId),
		State:          gitlab.String("opened"),
		WIP:            getWIPQueryParamPointer(shouldIncludeDrafts),
	})
	if err != nil {
		log.Printf("Failed to get merge requests for %v: %v\n", userId, err)
		return nil, err
	}

	return userMrs, nil
}

// FetchProjectMergeRequests fetches merge requests for a project from GitLab.
func FetchProjectMergeRequests(glabClient *glab.TGitlabClient, projectId string, usernames []string, shouldIncludeDrafts *bool) ([]*gitlab.MergeRequest, error) {
	var projectMrs []*gitlab.MergeRequest

	for _, username := range usernames {
		userMrs, _, err := glabClient.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			AuthorUsername: gitlab.String(username),
			State:          gitlab.String("opened"),
			WIP:            getWIPQueryParamPointer(shouldIncludeDrafts),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get merge request for %s: %w", username, err)
		}

		projectMrs = append(projectMrs, userMrs...)
	}

	return projectMrs, nil
}

// PrintMergeRequests prints the details of the merge requests to the console.
func PrintMergeRequests(mrs []*gitlab.MergeRequest) {
	for _, mr := range mrs {
		fmt.Printf("@%s: %s\n", mr.Author.Username, mr.WebURL)
	}
}

// OpenMergeRequests opens the URLs of the merge requests in the user's default browser.
func OpenMergeRequests(mrs []*gitlab.MergeRequest) error {
	for _, mr := range mrs {
		if err := openURL(mr.WebURL); err != nil {
			return err
		}
	}
	return nil
}

func GetMergeRequestsApprovedByMe(glabClient *glab.TGitlabClient, groupId string, myId int, shouldIncludeDrafts *bool) ([]*gitlab.MergeRequest, error) {
	mrsApprovedByMe, _, err := glabClient.MergeRequests.ListGroupMergeRequests(groupId, &gitlab.ListGroupMergeRequestsOptions{
		ApprovedByIDs: gitlab.ApproverIDs([]int{myId}),
		State:         gitlab.String("opened"),
		WIP:           getWIPQueryParamPointer(shouldIncludeDrafts),
	})
	if err != nil {
		log.Printf("Failed to get merge requests approved by me: %v\n", err)
		return nil, err
	}

	return mrsApprovedByMe, nil
}
