package music

import (
	"os/exec"
)

const cacheDir = "music_cache/"

// func CheckCache()

func DownloadMP3(url string) *exec.Cmd {
//<<<<<<< HEAD
//	cmd := fmt.Sprintf("-f mp3 -birate 64 -path %s -id %s",
//		cacheDir, url)
//
//	return exec.Command("youtube-dl", cmd)
//=======
	cmd := exec.Command("bash", "download_song.sh", url)

	return cmd
}
