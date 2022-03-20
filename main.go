package main

import (
	"flag"
	"fmt"
	"godiscordspeechbot/bot"
	"godiscordspeechbot/cogs"
	"godiscordspeechbot/commands"
	"godiscordspeechbot/utils"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	configPath     string
	matthewBot     *bot.Bot
	commandHandler *commands.CommandHandler
	cogHandler     *cogs.CogHandler
	logger         *utils.Logger
)

func init() {
	flag.StringVar(&configPath, "c", "", "configPath")
	flag.Parse()
}

func main() {
	var err error

	matthewBot, err = bot.New(configPath)

	if err != nil {
		fmt.Println("Error reading config", err)
		return
	}

	err = matthewBot.Login()

	if err != nil {
		fmt.Println("Error starting bot", err)
		return
	}

	commandHandler = commands.NewCommandHandler()

	commands.LoadDirectoryToHandler(commandHandler)
	logger = utils.NewLogger("data/logs/log.txt")
	quit := make(chan bool, 1)

	matthewBot.Session.AddHandler(onMessageReceive)
	matthewBot.Session.AddHandler(processCommands)
	matthewBot.Session.AddHandler(processCogs)

	cogHandler = cogs.NewCogHandler()
	go cogHandler.Run(*matthewBot)

	err = matthewBot.Session.Open()
	go utils.StartLogger(logger, quit)

	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc         // If we receive anything, we exit
	quit <- true // Gracefully stop the logging coroutine

	fmt.Println("SIGINT received, gracefully shutting down...")
	_ = matthewBot.Session.Close()
}

func onMessageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	toLog := "[" + m.Timestamp.String() + "]" + " [" + m.Author.Username + "] " + m.Content
	fmt.Println(toLog)

	logger.Messages <- toLog

	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func processCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author
	b := s.State.User

	if user.ID == b.ID || user.Bot {
		return
	}

	content := m.Content

	fmt.Println(matthewBot.Prefix, content)

	if !strings.HasPrefix(content, matthewBot.Prefix) {
		return
	}

	fields := strings.Fields(content)
	args := fields[1:]
	name := strings.ToLower(fields[0])

	if name == "%help" {
		commandHandler.Help(matthewBot, m)
		return
	}

	cmd, ok := commandHandler.Get(name[1:])

	if !ok {
		fmt.Println("Command doesn't exist", name)
		return
	}

	fmt.Println("Processing command", name, "with arguments", args)

	c := *cmd

	c(matthewBot, m, args)
}

func processCogs(s *discordgo.Session, ctx *discordgo.MessageCreate) {
	user := ctx.Author
	b := s.State.User

	if user.ID == b.ID || user.Bot {
		return
	}

	content := ctx.Content

	fmt.Println(matthewBot.Prefix, content)

	if !strings.HasPrefix(content, matthewBot.Prefix) {
		return
	}

	fields := strings.Fields(content)
	args := fields[1:]
	name := strings.ToLower(fields[0])

	if name == matthewBot.Prefix+"cog" {
		cog := cogs.GetCog(name)

		interval, err := strconv.ParseInt(args[0], 10, 32)

		if err != nil {
			matthewBot.Say(ctx, "Are you sure you entered a number?", 5)
			log.Print(err)
		}

		cogs.RegisterCog(*cogHandler, ctx, cog, int(interval))
		return
	}
}
