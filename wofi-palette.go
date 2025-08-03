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
	filename := "config.json"
	settings, err := readConfig(filename)

	fmt.Println(settings)
	fmt.Println(err)

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

	option := runCommand(woficmd)

	cleanOption := strings.TrimSpace(option)

	var powercmd *exec.Cmd

	switch cleanOption {
	case "Shutdown":
		powercmd = exec.Command("systemctl", "shutdown")
	case "Reboot":
		powercmd = exec.Command("systemctl", "reboot")
	case "Suspend":
		powercmd = exec.Command("systemctl", "suspend")
	case "Logout":
		powercmd = exec.Command("loginctl", "terminate-user", os.Getenv("USER"))
	case "Lock":
		powercmd = exec.Command("hyprlock")
	}

	if powercmd != nil {
		runCommand(powercmd)
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
