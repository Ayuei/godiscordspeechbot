package botcommands

import (
	"../../bot"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Join(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	_, err := b.JoinUserVoiceChannel(ctx.Author.ID)

	if err != nil {
		fmt.Println("Error joining Voice Channel", err)
	}
}

