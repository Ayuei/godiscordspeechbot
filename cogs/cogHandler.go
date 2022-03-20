package cogs

import (
	"../bot"
	"time"
)

type (
	CogFunc func(*bot.Bot, map[string]string, time.Duration)

	CogHandler struct {
		cogs CogList
	}

	Cog struct {
		Interval  time.Duration
		Arguments map[string]string
		Cogfunc   CogFunc
	}

	CogList []Cog

	CommandHandler struct {
		cogs CogList
	}
)

func NewCogHandler() *CogHandler {
	return &CogHandler{
		CogList{},
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
}

func (h *CogHandler) Run(bot bot.Bot) {
	for _, c := range h.cogs {
		go c.Cogfunc(&bot, c.Arguments, c.Interval)
	}
}
