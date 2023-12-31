package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var help = `macglab

Automate gathering your work on gitlab.com to save time.

v4.3.0

Complete documentation is available at https://github.com/mjburtenshaw/macglab
`

var rootCmd = &cobra.Command{
	Version: "5.0.1",
	Use:   "macglab",
	Short: "macglab automates gathering your work on gitlab.com to save time.",
	Long: help,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(help)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
