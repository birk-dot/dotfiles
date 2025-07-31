package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// WaybarOutput defines the structure for Waybar's JSON output
type WaybarOutput struct {
	Text    string `json:"text"`
	Alt     string `json:"alt"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

func main() {
	// Check if playerctl is installed
	if !isCommandAvailable("playerctl") {
		output := WaybarOutput{
			Text:    "Error: playerctl not found",
			Tooltip: "Please install playerctl to use this module.",
			Class:   "error",
		}
		printJSON(output)
		os.Exit(1)
	}

	playerStatus, err := getPlayerctlOutput("status")
	if err != nil || (playerStatus != "Playing" && playerStatus != "Paused") {
		// Spotify is not playing or not active
		output := WaybarOutput{
			Text:    " Offline",
			Alt:     "Spotify",
			Tooltip: "Spotify is not playing or not running.",
			Class:   "offline",
		}
		printJSON(output)
		os.Exit(0) // Exit cleanly if Spotify is not active
	}

	artist, _ := getPlayerctlOutput("metadata", "artist")
	title, _ := getPlayerctlOutput("metadata", "title")
	volumeStr, _ := getPlayerctlOutput("volume")

	// Convert volume to percentage
	volume, err := strconv.ParseFloat(volumeStr, 64)
	volumePercent := 0
	if err == nil {
		volumePercent = int(volume * 100)
	}

	icon := "" // Default to pause icon
	class := "paused"
	if playerStatus == "Playing" {
		icon = "" // Play icon
		class = "playing"
	}

	displayText := fmt.Sprintf("%s %s - %s", icon, artist, title)
	tooltipText := fmt.Sprintf("<b>Artist:</b> %s\n<b>Title:</b> %s\n<b>Status:</b> %s\n<b>Volume:</b> %d%%", artist, title, playerStatus, volumePercent)

	output := WaybarOutput{
		Text:    displayText,
		Alt:     "Spotify",
		Tooltip: tooltipText,
		Class:   class,
	}

	printJSON(output)
}

// isCommandAvailable checks if a given command exists in the system's PATH.
func isCommandAvailable(name string) bool {
	cmd := exec.Command("which", name)
	err := cmd.Run()
	return err == nil
}

// getPlayerctlOutput executes a playerctl command and returns its trimmed output.
// It takes the command and any arguments as strings.
func getPlayerctlOutput(args ...string) (string, error) {
	cmdArgs := []string{"--player=spotify"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("playerctl", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// printJSON marshals the WaybarOutput struct to JSON and prints it to stdout.
func printJSON(data WaybarOutput) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		// Fallback to plain text error if JSON marshalling fails
		fmt.Printf("{\"text\": \"JSON Error\", \"tooltip\": \"%v\"}", err)
		return
	}
	fmt.Println(string(jsonData))
}
