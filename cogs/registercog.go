package cogs

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func RegisterCog(handler CogHandler, ctx *discordgo.MessageCreate, cogFunc CogFunc, interval ...int) {
	duration := time.Duration(interval[0]) * time.Second

	arguments := make(map[string]string)
	arguments["server"] = ctx.ChannelID

	handler.RegisterCog(Cog{
		Interval: duration,
		ctx:      ctx,
		Cogfunc:  cogFunc,
	})
}
