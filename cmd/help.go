package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(helpCmd)
}

var AvailableCommands = `

macglab

Automate gathering your work on gitlab.com to save time.

v4.2.1

Commands:
- help: Prints help about a given command or a list of available commands if none provided.
- init: Initializes macglab.
- list: Prints GitLab Merge Request (MRs) authors and URLs to the terminal.
`

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Prints help",
	Long:  `Prints help for a given command, otherwise the help for the program.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(AvailableCommands)
	},
}
