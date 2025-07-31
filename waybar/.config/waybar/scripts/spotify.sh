#!/bin/bash

# This script fetches Spotify playback information and outputs it in JSON format
# for Waybar, allowing for custom display and control.

# Check if playerctl is installed
if ! command -v playerctl &> /dev/null; then
    echo '{"text": "Error: playerctl not found", "tooltip": "Please install playerctl to use this module."}'
    exit 1
fi

# Check if Spotify is running and active
# The 'playerctl --player=spotify status' command will return an error if Spotify is not running,
# or if it's running but not playing/paused (e.g., just opened).
# We check for both 'Playing' and 'Paused' states to ensure it's an active player.
player_status=$(playerctl --player=spotify status 2>/dev/null)

if [ "$player_status" = "Playing" ] || [ "$player_status" = "Paused" ]; then
    # Get metadata: artist and title
    artist=$(playerctl --player=spotify metadata artist 2>/dev/null)
    title=$(playerctl --player=spotify metadata title 2>/dev/null)

    # Get current volume
    volume=$(playerctl --player=spotify volume 2>/dev/null)
    # Convert volume to percentage and round to nearest integer
    volume_percent=$(printf "%.0f" "$(echo "$volume * 100" | bc)")

    # Determine icon based on playback status
    # You might need a Nerd Font installed for these icons to display correctly.
    if [ "$player_status" = "Playing" ]; then
        icon="" # Play icon (Font Awesome Solid Play)
        class="playing"
    else
        icon="" # Pause icon (Font Awesome Solid Pause)
        class="paused"
    fi

    # Construct the text to display in Waybar
    display_text="$icon $artist - $title"

    # Construct the tooltip text
    tooltip_text="<b>Artist:</b> $artist\n<b>Title:</b> $title\n<b>Status:</b> $player_status\n<b>Volume:</b> $volume_percent%"

    # Output JSON for Waybar
    echo "{\"text\": \"$display_text\", \"alt\": \"Spotify\", \"tooltip\": \"$tooltip_text\", \"class\": \"$class\"}"
else
    # Spotify is not playing or not active
    echo '{"text": " Offline", "alt": "Spotify", "tooltip": "Spotify is not playing or not running.", "class": "offline"}'
fi

