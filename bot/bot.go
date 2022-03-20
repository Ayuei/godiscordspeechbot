package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

// Config struct for fields to map to Bot from config file
type Config struct {
	Prefix        string `json:"prefix"`
	ServiceURL    string `json:"service_url"`
	BotToken      string `json:"bot_token"`
	OwnerID       string `json:"owner_id"`
	UseSharding   bool   `json:"use_sharding"`
	ShardID       int    `json:"shard_id"`
	ShardCount    int    `json:"shard_count"`
	DefaultStatus string `json:"default_status"`
	RiotAPIKey    string `json:"riot_api_key"`
	LostArkURL    string `json:"lost_ark_base_url"`
}

// Bot struct to abstract
type Bot struct {
	Session  *discordgo.Session
	LoggedIn bool
	Prefix   string
	config   Config
	RiotAPI  *golio.Client
}

// Login logs in the bot
func (b *Bot) Login() error {

	discord, err := discordgo.New(b.config.BotToken)
	fmt.Println(b.config.BotToken)

	riotClient := golio.NewClient(b.config.RiotAPIKey,
		golio.WithRegion(api.RegionOceania))

	b.RiotAPI = riotClient
	b.Session = discord
	b.LoggedIn = true

	fmt.Println("Creating bot")

	return err
}

// New constructs a new Bot
func New(configPath string) (b *Bot, e error) {
	body, err := ioutil.ReadFile(configPath)

	if err != nil {
		fmt.Println("Error loading config", err)
	}

	var conf Config
	_ = json.Unmarshal(body, &conf)

	b = new(Bot)
	b.LoggedIn = false
	b.Prefix = conf.Prefix
	b.config = conf

	fmt.Println("Loaded configuration")
	fmt.Println("Prefix", b.Prefix)

	return b, err
}

func (b Bot) Say(ctx *discordgo.MessageCreate, message string, deleteAfter ...int) *discordgo.Message {
	delMessageMillis := -1

	msg, err := b.Session.ChannelMessageSend(ctx.ChannelID, message)

	if err != nil {
		fmt.Println("Error has occurred while sending message", err)
	}

	if len(deleteAfter) > 0 {
		delMessageMillis = deleteAfter[0]
	}

	if delMessageMillis > 0 {
		go func(t int, mID string) {
			time.Sleep(time.Duration(t) * time.Second)

			_ = b.Session.ChannelMessageDelete(ctx.ChannelID, mID)
		}(delMessageMillis, msg.ID)
	}

	return msg
}

func (b Bot) SendMsgChannel(channel string, message string, deleteAfter ...int) *discordgo.Message {
	delMessageMillis := -1

	msg, err := b.Session.ChannelMessageSend(channel, message)

	if err != nil {
		fmt.Println("Error has occurred while sending message", err)
	}

	if len(deleteAfter) > 0 {
		delMessageMillis = deleteAfter[0]
	}

	if delMessageMillis > 0 {
		go func(t int, mID string) {
			time.Sleep(time.Duration(t) * time.Second)

			_ = b.Session.ChannelMessageDelete(channel, mID)
		}(delMessageMillis, msg.ID)
	}

	return msg
}

func (b Bot) SayEmbed(ctx *discordgo.MessageCreate, message *discordgo.MessageEmbed, deleteAfter ...int) *discordgo.Message {
	delMessageMillis := -1

	msg, err := b.Session.ChannelMessageSendEmbed(ctx.ChannelID, message)

	if err != nil {
		fmt.Println("Error has occurred while sending message", err)
	}

	if len(deleteAfter) > 0 {
		delMessageMillis = deleteAfter[0]
	}

	if delMessageMillis > 0 {
		go func(t int, mID string) {
			time.Sleep(time.Duration(t) * time.Second)

			_ = b.Session.ChannelMessageDelete(ctx.ChannelID, mID)
		}(delMessageMillis, msg.ID)
	}

	return msg
}

func (b Bot) findUserVoiceState(userid string) (*discordgo.VoiceState, error) {
	for _, guild := range b.Session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("could not find users voice state")
}

func (b *Bot) JoinUserVoiceChannel(userID string) (*discordgo.VoiceConnection, error) {
	// Find a user's current voice channel
	vs, err := b.findUserVoiceState(userID)
	if err != nil {
		return nil, err
	}

	// Join the user's channel and start unmuted and deafened.
	return b.Session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, false)
}

func (b *Bot) PlaySong(songPath string, v *discordgo.VoiceConnection) chan bool {
	quit := make(chan bool)
	dgvoice.PlayAudioFile(v, songPath, quit)

	return quit
}

func (b *Bot) GetLostArkURL() string {
	return b.config.LostArkURL
}
