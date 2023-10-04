package main

import (
	"flag"

	"github.com/mjburtenshaw/macglab/mrs"
)

func main() {
	browserFlag := flag.Bool("browser", false, "Open merge requests in the browser")

	flag.Parse()

	mergeRequests := mrs.FetchMergeRequests()
	mrs.PrintMergeRequests(mergeRequests)

	if *browserFlag {
		mrs.OpenMergeRequests(mergeRequests)
	}
}
