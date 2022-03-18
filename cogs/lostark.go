package cogs

import (
	"godiscordspeechbot/bot"
	"time"
)

func LostArkCog(b *bot.Bot, args map[string]string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	channelID, ok := args["channelID"]

	if !ok {
		b.SendMsgChannel(channelID, "Something went wrong", 5)
		return
	}

	for {
		select {
		case <-ticker.C:

		}
	}
}
