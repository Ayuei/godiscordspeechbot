package commandUtils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"net/url"
)


type SpeechMessage struct{
	PCM []int16
}

const addr = "192.168.20.8:8080"

func Start(receive chan *discordgo.Packet, output chan string){

	u := url.URL{Scheme: "ws", Host: addr, Path: "/stream"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		select {
		case data := <- receive:
			err = conn.WriteJSON(SpeechMessage{PCM: data.PCM})
			err = conn.WriteJSON(SpeechMessage{PCM: []int16{}})

			if err != nil {
				fmt.Println(err)
				return
			}
			_, p, _ := conn.ReadMessage()

			fmt.Println(string(p))

			output <- string(p)
		}
	}

}
