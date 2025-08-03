package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Paths [][]string `json:"paths"`
}

type FormattedConfig struct {
	Name    string
	Command *exec.Cmd
}

func main() {
	homeFolder, homeFolderExists := os.LookupEnv("HOME")
	configFolder := homeFolder + "/.config/immaculate-wofi/"
	configFile := configFolder + "wofi-palette.json"

	if !homeFolderExists {
		log.Fatal("HOME environment variable not set")
	}

	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		err := os.MkdirAll(configFolder, 0755)
		if err != nil {
			log.Fatal("Error creating directory:", err)
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		warning := fmt.Sprintf(`{"paths": [["Add your paths in %s", "echo hi"]]}`, configFile)
		err := os.WriteFile(configFile, []byte(warning), 0755)

		if err != nil {
			log.Fatal("Error creating config file:", err)
		}
	}

	settings, err := readConfig(configFile)

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var options []FormattedConfig

	for _, path := range settings.Paths {
		if len(path) == 2 {
			name := path[0]
			cmd := path[1]

			// Create a new command with the provided command string
			command := exec.Command("sh", "-c", cmd)

			// Append the formatted option to the options slice
			options = append(options, FormattedConfig{
				Name:    name,
				Command: command,
			})
		} else {
			log.Printf("Skipping invalid path entry: %v", path)
		}
	}

	optionsString := ""

	for _, option := range options {
		name := option.Name
		optionsString += name + "\n"
	}

	woficmd := exec.Command(
		"wofi",
		"--show", "dmenu",
		"--prompt", "Command Palette",
	)

	woficmd.Stdin = bytes.NewBufferString(optionsString)

	selectedOption := strings.TrimSpace(runCommand(woficmd))

	for _, option := range options {
		if option.Name == selectedOption && option.Command != nil {
			runCommand(option.Command)
		}
	}

}

func runCommand(cmd *exec.Cmd) string {
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	return string(output)
}

func readConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var settings Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		return Config{}, err
	}

	return settings, nil
}
