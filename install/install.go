package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("macglab: 🏗️  Installing macglab...")
	homeUri := os.Getenv("HOME")
	if homeUri == "" {
		log.Fatal("macglab: 🏚️ Couldn't find HOME environment variable.")
	}

    macglabUri := fmt.Sprintf("%s/.macglab", homeUri)

	isNewInstall := false
    if _, err := os.Stat(macglabUri); os.IsNotExist(err) {
			fmt.Println("macglab: 🆕 No previous installation detected. *cracks knuckles* Starting from scratch...")
					isNewInstall = true

			fmt.Println("macglab: 🏠 Making home directory for macglab...")
			cmd := exec.Command("mkdir", macglabUri)
			err := cmd.Run()
			if err != nil {
				log.Fatal("macglab: 💀 Couldn't create macglab config directory.")
			}

			fmt.Println("macglab: 🐚 Adding environment variables...")
			shConfigUrl := fmt.Sprintf("%s/.zshrc", homeUri)
			shConfig, err := os.OpenFile(shConfigUrl, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal("macglab: 💀 Couldn't write to shell config file.")
			}
			defer shConfig.Close()

			writer := bufio.NewWriter(shConfig)

			lines := []string{
				"",
				"# [`macglab`](https://github.com/mjburtenshaw/macglab)",
				`export MACGLAB="${HOME}/.macglab"`,
				`export PATH="${MACGLAB}/bin:${PATH}"`,
				`alias macglab="${MACGLAB}/bin"`,
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

	fmt.Println("macglab: 🔨 Building executable...")
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("💀 Error getting current working directory: %v", err)
	}

    cmd := exec.Command("go", "build", "-o", "bin", "main.go")
    cmd.Dir = cwd
    err = cmd.Run()
    if err != nil {
        log.Fatalf("💀 Error building: %v", err)
    }

	fmt.Println("macglab: 🚚 Moving executable to macglab home...")
	binPath := fmt.Sprintf("%s/bin", macglabUri)
	cmd = exec.Command("mv", "bin", binPath)
    err = cmd.Run()
    if err != nil {
        log.Fatalf("💀 Error moving executable: %v", err)
    }

	macglabConfigUrl := fmt.Sprintf("%s/config.yml", macglabUri)
	if isNewInstall {
		fmt.Println("macglab: 📜 Making a new config file...")
		cmd = exec.Command("cp", "config.sample.yml", macglabConfigUrl)
		err = cmd.Run()
		if err != nil {
			log.Fatal("macglab: 💀 Couldn't create macglab config file.")
		}
	}

	fmt.Println("macglab: 🎉 Successfully installed!")
	
	if isNewInstall {
		fmt.Printf("macglab: 📜 Created a new config file at %s. Please open it and define values.\n", macglabConfigUrl)
		fmt.Println("macglab: 🐚 Re-source your shell session or open a new terminal, then run `macglab` and watch the magic happen!")
	}
}
