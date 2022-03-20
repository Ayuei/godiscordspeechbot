package botcommands

import (
	"bioreddit_indexer/utils"
	"github.com/bwmarrin/discordgo"
	"godiscordspeechbot/bot"
	"godiscordspeechbot/cogs"
	"time"
)

func RegisterCog(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	cog := cogs.GetCog(args[2])
	duration, err := time.ParseDuration(args[0])
	utils.CheckError(err, "Duration")
	arguments := make(map[string]string)
	arguments["server"] = ctx.ChannelID

	b.CogHandler.RegisterCog(cogs.Cog{
		Interval: duration,
		Arguments: arguments,
		Cogfunc: cog,
	})
}
