#!/bin/sh

youtube-dl -x --audio-format opus --audio-quality 64k -o "/home/dietpi/github/golang_discord_assistant/src/music_cache/%(id)s.%(ext)s" "${1}"
