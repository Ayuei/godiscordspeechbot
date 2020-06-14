package botcommands

import (
	"../../bot"
	"./commandUtils"
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func Speak(b *bot.Bot, ctx *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		b.Say(ctx, "You need to tell me what to say!", 3)
	}

	v, found := b.Session.VoiceConnections[ctx.GuildID]

	if !found {
		b.Say(ctx, "You need to me !join first!", 3)
		return
	}

	lang := "en"
	lastArg := args[len(args)-1]

	if len(lastArg) == 2 && !strings.HasSuffix(lastArg, ".") {
		lang = lastArg
		args[len(args)-1] = ""
	}

	text := strings.Join(args, " ")

	log.Println("Language", lang)
	speech := commandUtils.Speech{
		Folder:   "mp3",
		Language: lang,
	}

	fp, err := speech.Speak(text)

	if err != nil {
		fmt.Println("Fatal error", err)
	}

	go dgvoice.PlayAudioFile(v, fp, make(chan bool))
}
