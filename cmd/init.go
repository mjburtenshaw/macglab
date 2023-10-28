package cmd

import (
	"log"
	"os"

	"github.com/mjburtenshaw/macglab/config"
	"github.com/mjburtenshaw/macglab/files"
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
		log.Println("macglab: installing macglab...")

		isNewInstall := false
		if _, err := os.Stat(files.MacglabUri); os.IsNotExist(err) {
			log.Println("macglab: no previous installation detected. *cracks knuckles* Starting from scratch...")
			isNewInstall = true

			log.Println("macglab: demanding a home directory for macglab...")
			if err := files.DemandDir(files.MacglabUri); err != nil {
				log.Fatalf("macglab: couldn't create macglab config directory: %s", err)
			}

			log.Println("macglab: adding environment variables...")
			if err := config.AddEnv(files.ShConfigUrl); err != nil {
				log.Fatalf("macglab: couldn't add environment variables: %s", err)
			}
		}

		if isNewInstall {
			log.Println("macglab: making a new config file...")
            if err := config.Create(files.SampleConfigUrl, files.MacglabConfigUrl); err != nil {
                log.Fatalf("macglab: couldn't add config: %s", err)
            }
		}

		log.Println("macglab: 🎉 Successfully installed!")

		if isNewInstall {
			log.Printf("macglab: created a new config file at %s. Please open it and define values.\n", files.MacglabConfigUrl)
			log.Println("macglab: re-source your shell session or open a new terminal, then run `macglab list` and watch the magic happen!")
		}
	},
}
