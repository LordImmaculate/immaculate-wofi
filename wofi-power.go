package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	options := []string{
		"Shutdown",
		"Reboot",
		"Suspend",
		"Logout",
		"Lock",
	}

	optionsString := ""

	for _, option := range options {
		optionsString += option + "\n"
	}

	woficmd := exec.Command(
		"wofi",
		"--show", "dmenu",
		"--prompt", "Power options",
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
