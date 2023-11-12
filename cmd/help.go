package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(helpCmd)
	helpCmd.AddCommand(initHelpCmd)
}

var AvailableCommands = `macglab

Automate gathering your work on gitlab.com to save time.

v4.2.1

Commands:
- help: Prints help about a given command or a list of available commands if none provided.
- init: Initializes macglab.
- list: Prints GitLab Merge Request (MRs) authors and URLs to the terminal.
`

var InitHelp = `init

Initializes macglab.

init does the following:

1. Checks if there's a previous installation. Exits if so.
2. Demands a home directory for this program on your machine.
3. Adds required environment variables to your shell config file.
4. Makes a new config file.

The config directory is created at ~/.macglab.

The config file is located at ~/.macglab/config.yml

We support the following shells:
- zsh.

Options:

-h, --help: Print these docs.
`

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Prints help",
	Long:  "Prints help for a given command, otherwise the help for the program.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(AvailableCommands)
	},
}

var initHelpCmd = &cobra.Command{
	Use:   "init",
	Short: "Prints help for the init command",
	Long:  "Prints help for the init command",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(InitHelp)
	},
}
