package botcommands

import (
	"../../bot"
	"../../cogs"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func RegisterCog(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	cog := cogs.GetCog(args[2])
	duration, err := time.ParseDuration(args[0])

	if err != nil {
		log.Fatal(err)
	}

	arguments := make(map[string]string)
	arguments["server"] = ctx.ChannelID

	b.CogHandler.RegisterCog(cogs.Cog{
		Interval: duration,
		Arguments: arguments,
		Cogfunc: cog,
	})
}
