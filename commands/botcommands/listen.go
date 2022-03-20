package botcommands

import (
	"godiscordspeechbot/bot"
	"godiscordspeechbot/commands/botcommands/commandUtils"

	"github.com/bwmarrin/discordgo"
)

func Listen(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	v, ok := b.Session.VoiceConnections[ctx.GuildID]

	if ok != true {
		_ = b.Say(ctx, "You need to make me !join first", 5)
		return
	}

	recv := make(chan *discordgo.Packet, 2)
	send := make(chan string, 2)

	go commandUtils.ReceiveAndConvertPCM(v, recv)
	go commandUtils.Start(recv, send)

	//v.Speaking(true)
	//defer v.Speaking(false)

	for {
		b.Say(ctx, <-send)

		if !ok {
			return
		}
	}
}
