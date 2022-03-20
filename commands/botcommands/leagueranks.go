package botcommands

import (
	"fmt"
	"godiscordspeechbot/bot"
	"strings"

	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/bwmarrin/discordgo"
)

func formatRankedMessage(players string, ranks string, winRates string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "",
		Color: 0xffa500,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Player Names",
				Value:  players,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Ranks",
				Value:  ranks,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "W/R",
				Value:  winRates,
				Inline: true,
			},
		},
	}
}

func LookupGame(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	name := strings.Join(args, " ")
	summoner, err := b.RiotAPI.Riot.Summoner.GetByName(name)

	if err != nil {
		b.Say(ctx, "This person doesn't exist?!", 3)
		fmt.Println("Error looking up summoner", err)
		return
	}

	matchData, err := b.RiotAPI.Riot.Spectator.GetCurrent(summoner.ID)

	if err != nil {
		b.Say(ctx, "This person is not in a game!", 3)
		fmt.Println("Error looking up match", err)
		return
	}

	var p *lol.CurrentGameParticipant // Provide a type hint
	results := make(chan *lol.LeagueItem)

	for _, p = range matchData.Participants {
		go func(cgp *lol.CurrentGameParticipant) {
			queues, _ := b.RiotAPI.Riot.League.ListBySummoner(cgp.SummonerID)
			var q *lol.LeagueItem

			flag := true

			for _, q = range queues {
				if q.QueueType == string(lol.QueueRankedSolo) {
					results <- q
					flag = false
				}
			}

			if flag {
				results <- &lol.LeagueItem{
					SummonerName: cgp.SummonerName,
					Tier:         "",
					Rank:         "Unranked",
				}
			}
		}(p)
	}

	var names strings.Builder
	var ranks strings.Builder
	var winRates strings.Builder

	tasksDone := 0

	for res := range results {
		names.WriteString(res.SummonerName + "\n")
		if res.Rank != "Unranked" {
			ranks.WriteString(res.Tier + " " + res.Rank + "\n")
			winRate := float32(res.Wins) / float32(res.Wins+res.Losses)
			winRateString := fmt.Sprintf("%.2f\n", winRate)
			winRates.WriteString(winRateString)
		} else {
			ranks.WriteString(res.Rank + "\n")
			winRates.WriteString("N/A\n")
		}

		tasksDone++

		if tasksDone == 10 {
			close(results)
			break
		}
	}

	b.SayEmbed(ctx, formatRankedMessage(names.String(),
		ranks.String(),
		winRates.String(),
	))
}
