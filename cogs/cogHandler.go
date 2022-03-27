package cogs

import (
	"github.com/bwmarrin/discordgo"
	"godiscordspeechbot/bot"
	"log"
	"time"
)

type (
	CogFunc func(b *bot.Bot, context *discordgo.MessageCreate, interval time.Duration)

	CogHandler struct {
		cogs     CogList
		cogsChan chan Cog
	}

	Cog struct {
		Interval time.Duration
		ctx      *discordgo.MessageCreate
		Cogfunc  CogFunc
	}

	CogList []Cog
)

func NewCogHandler() *CogHandler {
	return &CogHandler{
		CogList{},
		make(chan Cog),
	}
}

func (h CogHandler) GetCogs() CogList {
	// Return map attribute
	return h.cogs
}

func (h CogHandler) Get(idx int) (*Cog, bool) {
	if idx > len(h.cogs) {
		return nil, false
	}

	cog := h.cogs[idx]

	return &cog, true
}

func (h *CogHandler) RegisterCog(cog Cog) {
	h.cogs = append(h.cogs, cog)
	h.cogsChan <- cog

}

func (h *CogHandler) Run(bot *bot.Bot) {
	log.Println("Running Cog Handler")

	for {
		select {
		case c := <-h.cogsChan:
			log.Println("Cog Handler")
			go c.Cogfunc(bot, c.ctx, c.Interval)
		}
	}
}
