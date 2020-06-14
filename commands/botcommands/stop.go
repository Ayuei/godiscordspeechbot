package botcommands

import (
	"../../bot"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Stop(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	v, ok := b.Session.VoiceConnections[ctx.GuildID]

	if ok != true {
		_ = b.Say(ctx, "m9, I ain't even there", 3)
		return
	}

	fmt.Println("Closing voice connection")

	v.Disconnect()
}

