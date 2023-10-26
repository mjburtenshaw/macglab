package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes macglab",
	Long:  `Initializes macglab`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("macglab: 🏗️  Installing macglab...")

		isNewInstall := false
		if _, err := os.Stat(config.MacglabUri); os.IsNotExist(err) {
			fmt.Println("macglab: 🆕 No previous installation detected. *cracks knuckles* Starting from scratch...")
			isNewInstall = true

			if err := config.DemandConfigDir(); err != nil {
				log.Fatalf("macglab: couldn't create macglab config directory: %s", err)
			}

			if err := config.AddEnv(config.ShConfigUrl); err != nil {
				log.Fatalf("macglab: couldn't add environment variables: %s", err)
			}
		}

		if isNewInstall {
			fmt.Println("macglab: 📜 Making a new config file...")
			cmd := exec.Command("cp", "config.sample.yml", config.MacglabConfigUrl)
			err := cmd.Run()
			if err != nil {
				log.Fatal("macglab: 💀 Couldn't create macglab config file.")
			}
		}

		fmt.Println("macglab: 🎉 Successfully installed!")

		if isNewInstall {
			fmt.Printf("macglab: 📜 Created a new config file at %s. Please open it and define values.\n", config.MacglabConfigUrl)
			fmt.Println("macglab: 🐚 Re-source your shell session or open a new terminal, then run `macglab list` and watch the magic happen!")
		}
	},
}
