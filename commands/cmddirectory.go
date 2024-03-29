package commands

import (
	"godiscordspeechbot/commands/botcommands"
)

var lookUpCommand = CmdInfo{command: botcommands.LookupGame, help: "Look up the ranked statistics of the players in your game"}
var playCommand = CmdInfo{command: botcommands.Play, help: "Play a Youtube Video"}
var echoCommand = CmdInfo{command: botcommands.Echo, help: "Jason 2.0"}
var joinCommand = CmdInfo{command: botcommands.Join, help: "Summons me to your voice channel"}
var stopCommand = CmdInfo{command: botcommands.Stop, help: "Desummons me"}
var listenCommand = CmdInfo{command: botcommands.Listen, help: "I will listen and transcribe"}
var speakCommand = CmdInfo{command: botcommands.Speak, help: "I will say what you tell me to"}
var serverStatus = CmdInfo{command: botcommands.ServerStatus, help: "I will say what you tell me to"}

func LoadDirectoryToHandler(h *CommandHandler) {
	h.RegisterCommand("league_ranks", lookUpCommand)
	h.RegisterCommand("play", playCommand)
	h.RegisterCommand("join", joinCommand)
	h.RegisterCommand("echo", echoCommand)
	h.RegisterCommand("stop", stopCommand)
	h.RegisterCommand("listen", listenCommand)
	h.RegisterCommand("speak", speakCommand)
	h.RegisterCommand("isup", serverStatus)
}
