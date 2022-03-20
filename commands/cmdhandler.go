package commands

import (
	"github.com/bwmarrin/discordgo"
	"godiscordspeechbot/bot"
	"strings"
)

// Code refactored and commented from
// https://github.com/ducc/GoMusicBot/blob/master/src/framework/command.go

type (
	Command func(*bot.Bot, *discordgo.MessageCreate, []string)

	CmdInfo struct {
		command Command
		help    string
	}

	// We never need to change the struct
	// So we by pass by value rather than reference
	CmdMap map[string]CmdInfo

	CommandHandler struct {
		cmds CmdMap
	}
)

func NewCommandHandler() *CommandHandler {
	// Return a command handler with an zero initialized map
	return &CommandHandler{
		make(CmdMap),
	}
}

func (h CommandHandler) GetCmds() CmdMap {
	// Return map attribute
	return h.cmds
}

func (h CommandHandler) Get(name string) (*Command, bool) {
	cmd, ok := h.cmds[name]

	return &cmd.command, ok
}

func (h *CommandHandler) RegisterCommand(name string, cmd CmdInfo) {
	h.cmds[name] = cmd

	// ! Create a shortcut alias for our command
	if len(name) > 1 {
		h.cmds[name[:1]] = cmd
	}
}

func (c CmdInfo) GetHelp() string {
	return c.help
}

func formatHelpCommand(commands string, helpString string) *discordgo.MessageEmbed {

	return &discordgo.MessageEmbed{
		Title: "List of available commands",
		Color: 0x6a0dad,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Did someone say secret commands",
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Command",
				Value:  commands,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Help",
				Value:  helpString,
				Inline: true,
			},
		},
	}
}

func (h *CommandHandler) Help(b *bot.Bot, ctx *discordgo.MessageCreate) {

	args := strings.Split(ctx.Content, " ")
	arg := "all"
	var cmds map[string]CmdInfo

	if len(args) > 1 {
		arg = args[1]
	}

	if arg == "all" {
		cmds = h.GetCmds()
	} else {
		cmds = make(map[string]CmdInfo)
		cmd, found := h.GetCmds()[arg]

		if !found {
			b.Say(ctx, "Command doesn't exist!", 3)
		}

		cmds[arg] = cmd
	}

	var commands strings.Builder
	var helpStrings strings.Builder
	var name string
	var cmd CmdInfo

	for name, cmd = range cmds {
		if len(name) > 1 {
			commands.WriteString(name + "\n")
			helpStrings.WriteString(cmd.GetHelp() + "\n")
		}
	}

	msg := formatHelpCommand(commands.String(), helpStrings.String())

	b.SayEmbed(ctx, msg)
}
