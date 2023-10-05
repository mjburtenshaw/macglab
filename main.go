package main

import (
	"flag"

	"github.com/mjburtenshaw/macglab/mrs"
)

func main() {
	browserFlag := flag.Bool("browser", false, "Open merge requests in the browser")

	flag.Parse()

	groupMrs := mrs.FetchGroupMergeRequests()
	mrs.PrintMergeRequests(groupMrs)

	if *browserFlag {
		mrs.OpenMergeRequests(groupMrs)
	}
}
