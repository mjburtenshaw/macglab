package cmd

import (
	"bufio"
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

		goPath := os.Getenv("GOPATH")

		if goPath == "" {
			log.Fatal("macglab: 🏚️ Couldn't find GOPATH environment variable.")
		}

		isNewInstall := false
		if _, err := os.Stat(config.MacglabUri); os.IsNotExist(err) {
			fmt.Println("macglab: 🆕 No previous installation detected. *cracks knuckles* Starting from scratch...")
			isNewInstall = true

			fmt.Println("macglab: 🏠 Making home directory for macglab...")
			cmd := exec.Command("mkdir", config.MacglabUri)
			err := cmd.Run()
			if err != nil {
				log.Fatal("macglab: 💀 Couldn't create macglab config directory.")
			}

			fmt.Println("macglab: 🐚 Adding environment variables...")
			shConfig, err := os.OpenFile(config.ShConfigUrl, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal("macglab: 💀 Couldn't write to shell config file.")
			}
			defer shConfig.Close()

			writer := bufio.NewWriter(shConfig)

			lines := []string{
				"",
				"# [`macglab`](https://github.com/mjburtenshaw/macglab)",
				"",
				`export MACGLAB="${HOME}/.macglab"`,
				`export PATH="${GOPATH}/bin/macglab:${PATH}"`,
				"",
			}

			for _, line := range lines {
				_, err := writer.WriteString(line + "\n")
				if err != nil {
					log.Fatal("macglab: 💀 Couldn't write line to shell config file.")
				}
			}

			writer.Flush()
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
