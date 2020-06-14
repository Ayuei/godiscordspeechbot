package botcommands

import (
	"../../bot"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func Echo(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	v, ok := b.Session.VoiceConnections[ctx.GuildID]

	if ok != true {
		_ = b.Say(ctx, "You need to make me !join first", 5)
		return
	}

	recv := make(chan *discordgo.Packet, 2)
	send := make(chan []int16, 2)

	go dgvoice.ReceivePCM(v, recv)
	go dgvoice.SendPCM(v, send)

	v.Speaking(true)
	defer v.Speaking(false)

	for {
		p, ok := <-recv

		if !ok {
			return
		}

		send <- p.PCM
	}
}

