package botcommands

import (
	"encoding/json"
	"godiscordspeechbot/bot"
	"godiscordspeechbot/utils"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Status struct {
	Status int `json:"status"`
	Data   struct {
		Valtan string `json:"Valtan"`
	} `json:"data"`
}

func ServerStatus(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {

	log.Println("Getting server status.")

	resp := utils.CurlGet(b.GetLostArkURL(), "/server/Valtan")

	var status Status

	err := json.Unmarshal(resp, &status)

	if err != nil {
		b.Say(ctx, "Unable to get server status.", 5)
	}

	b.Say(ctx, "Server is "+status.Data.Valtan)
}
