package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "macglab",
  Short: "macglab automates gathering your work on gitlab.com to save time.",
  Long: `macglab automates gathering your work on gitlab.com to save time.

					This program lists all GitLab Merge Requests (MRs) based on:
					
					- Open state
					- Specified usernames and/or projects
					- Specified group
					
					Complete documentation is available at https://github.com/mjburtenshaw/macglab`,
  Run: func(cmd *cobra.Command, args []string) {
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
