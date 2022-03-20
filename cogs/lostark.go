package cogs

import (
	"../bot"
	"../utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var NewsCategories = []string{"updates", "events", "release-notes", "general"}

type NewsResponse struct {
	Status int `json:"status"`
	Data   []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnail   string `json:"thumbnail"`
		URL         string `json:"url"`
		PublishDate string `json:"publishDate"`
		Excerpt     string `json:"excerpt"`
	} `json:"data"`
}

type ForumResponse struct {
	Status int `json:"status"`
	Data   []struct {
		Title     string    `json:"title"`
		PostBody  string    `json:"post_body"`
		CreatedAt time.Time `json:"created_at"`
		URL       string    `json:"url"`
		Author    string    `json:"author"`
	} `json:"data"`
}

type DiscordMsg struct {
	Type  string
	Title string
	Body  string
	URL   string
}

func curlGet(baseURL string, path string) string {
	baseURL = strings.Trim(baseURL, "/")
	req, err := http.NewRequest("GET", baseURL+"/"+path, nil)

	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Unable to retrieve news:", err.Error())
	}

	defer resp.Body.Close()

	if resp != nil && resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)
		log.Println(bodyString)

		return bodyString
	}

	return ""
}

func GetNews(b *bot.Bot, category string) NewsResponse {
	respString := curlGet(b.GetLostArkURL(), "/news/"+category)

	var news NewsResponse

	err := json.Decoder.Decode(&respString, news)

	if err != nil {
		log.Fatal(err)
	}

	return news
}

func GetForumUpdates(b *bot.Bot) ForumResponse {
	respString := curlGet(b.GetLostArkURL(), "/v1/forums")

	var frm ForumResponse
	err := json.Decoder.Decode(&respString, frm)

	if err != nil {
		log.Fatal(err)

	}

	return frm
}

func GetNewsItems(hashmap map[uint32]bool, response NewsResponse) []DiscordMsg {
	var msgs []DiscordMsg

	for _, datum := range response.Data {
		hash := utils.Hash(datum.URL + datum.Title)

		// If it's not in the hashmap, it's new
		if _, ok := hashmap[hash]; !ok {
			hashmap[hash] = true
			msgs = append(msgs, DiscordMsg{
				"News",
				datum.Title,
				datum.Description,
				datum.URL,
			})
		} else {
			// Since these will appear in chronological order, we break on first
			break
		}
	}

	return msgs
}

func GetForumItems(hashmap map[uint32]bool, response ForumResponse) []DiscordMsg {
	var msgs []DiscordMsg

	for _, datum := range response.Data {
		hash := utils.Hash(datum.URL + datum.Title)

		// If it's not in the hashmap, it's new
		if _, ok := hashmap[hash]; !ok {
			hashmap[hash] = true
			msgs = append(msgs, DiscordMsg{
				"Forum",
				datum.Title,
				datum.PostBody,
				datum.URL,
			})
		} else {
			// Since these will appear in chronological order, we break on first
			break
		}
	}
	return msgs
}

func LostArkCog(b *bot.Bot, args map[string]string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	channelID, ok := args["channelID"]
	newsHashMap := make(map[uint32]bool)
	forumHashap := make(map[uint32]bool)

	if !ok {
		b.SendMsgChannel(channelID, "Something went wrong", 5)
		return
	}

	for {
		select {
		case <-ticker.C:
			for _, category := range NewsCategories {
				news := GetNewsItems(newsHashMap, GetNews(b, category))

				if len(news) > 0 {
					for _, newsItem := range news {
						b.SendMsgChannel(channelID, newsItem.URL)
					}
				}
			}

			forumUpdates := GetForumItems(forumHashap, GetForumUpdates(b))

			if len(forumUpdates) > 0 {
				for _, forumUpdate := range forumUpdates {
					b.SendMsgChannel(channelID, forumUpdate.URL)
				}
			}
		}
	}
}
