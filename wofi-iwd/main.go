package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	adapter := "wlan0"

	scancmd := exec.Command("iwctl", "station", adapter, "scan")
	runCommand(scancmd)

	networkscmd := exec.Command("iwctl", "station", adapter, "get-networks")
	output := runCommand(networkscmd)

	re := regexp.MustCompile("\x1B\\[[0-9;]*m")
	output = re.ReplaceAllString(output, "")
	lines := strings.Split(output, "\n")

	if len(lines) > 4 {
		lines = lines[4:]
	} else {
		errorcmd := exec.Command("notify-send", "No networks avaible.")
		runCommand(errorcmd)
		log.Fatal("No networks available.")
	}

	networkList := []string{}
	for _, line := range lines {
		if len(line) > 0 {
			networkList = append(networkList, line)
		}
	}

	woficmd := exec.Command("wofi", "--show dmenu", "--prompt 'Select a network'")
	woficmd.Stdin = bytes.NewBufferString(strings.Join(networkList, "\n"))
	selectedNetwork := runCommand(woficmd)
	fmt.Print(selectedNetwork)
}

func runCommand(cmd *exec.Cmd) string {
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
		log.Fatalf("Command failed: %v", err)
	}

	return string(output)
}
