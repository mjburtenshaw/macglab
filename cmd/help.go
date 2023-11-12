package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(helpCmd)
}

var MacglabHelp = `macglab

`

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Prints help",
	Long:  `Prints help for a given command, otherwise the help for the program.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print(``)
	},
}
