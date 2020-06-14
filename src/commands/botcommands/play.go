package botcommands

import (
	"../../bot"
	"../../music"
	"fmt"
	"strings"
	"os"
	"github.com/bwmarrin/discordgo"
)

const basePath = "/home/dietpi/github/golang_discord_assistant/src/music_cache/"

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func Play(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
//<<<<<<< HEAD
//
//	if len(args) == 0 {
//		b.Say(ctx, "You need to specify a link!", 3)
//	}
//
//	url := args[0]
//	cmd := music.DownloadMP3(url)
//=======
	url := args[0]
	fp := basePath+strings.Split(url, "=")[1]+".opus"

	ok, err := exists(fp)
//>>>>>>> a9bea77be79a8b8d31b6657ce6f08d7c6c3183e4

	if !ok || err != nil {
		fmt.Println("Downloading new Mp3")
		cmd := music.DownloadMP3(url)
		err := cmd.Run()

		if err != nil {
			fmt.Println("Fatal error", err)
		}
	}

	v, found := b.Session.VoiceConnections[ctx.GuildID]

	if !found {
		b.Say(ctx, "You need to me !join first!", 3)
		return
	}

	b.PlaySong(fp, v)
}
