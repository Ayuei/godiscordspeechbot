package music

import (
	"../bot"
	"github.com/bwmarrin/discordgo"
)

type SongQueue struct {
	list    []Song
	current *Song
	Running bool
}

func (queue SongQueue) Get() []Song {
	return queue.list
}

func (queue *SongQueue) Set(list []Song) {
	queue.list = list
}

func (queue *SongQueue) Add(song Song) {
	queue.list = append(queue.list, song)
}

func (queue SongQueue) HasNext() bool {
	return len(queue.list) > 0
}

func (queue *SongQueue) Next() Song {
	song := queue.list[0]
	queue.list = queue.list[1:]
	queue.current = &song
	return song
}

func (queue *SongQueue) Clear() {
	queue.list = make([]Song, 0)
	queue.Running = false
	queue.current = nil
}

func (queue *SongQueue) Start(b bot.Bot, ctx *discordgo.MessageCreate) {
	queue.Running = true
	for queue.HasNext() && queue.Running {
		song := queue.Next()
		b.Say(ctx, "Now playing `" + song.Title + "`.")
		v, found := b.Session.VoiceConnections[ctx.GuildID]

		if !found {
			b.Say(ctx, "You need to me !join first!", 3)
			return
		}

		b.PlaySong(song.Media, v)
	}
	if !queue.Running {
		b.Say(ctx, "Stopped playing.")
	} else {
		b.Say(ctx, "Finished queue.")
	}
}

func (queue *SongQueue) Current() *Song {
	return queue.current
}

func (queue *SongQueue) Pause() {
	queue.Running = false
}

func newSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]Song, 0)
	return queue
}