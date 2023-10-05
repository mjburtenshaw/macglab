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

func FetchGroupMergeRequests() [] *gitlab.MergeRequest {
	var groupMrs [] *gitlab.MergeRequest
	
	for _, username := range config.Usernames {
		userMrs, _, err := glab.Client.MergeRequests.ListGroupMergeRequests(config.GroupId, &gitlab.ListGroupMergeRequestsOptions{
			State: gitlab.String("opened"),
			AuthorUsername: gitlab.String(username),
		})
		if err != nil {
			log.Fatalf("💀 Failed to get merge request for %s: %v", username, err)
		}

		groupMrs = append(groupMrs, userMrs...)
	}
	
	return groupMrs
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
