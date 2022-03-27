package cogs

import (
	"encoding/json"
	"github.com/abadojack/whatlanggo"
	"github.com/bwmarrin/discordgo"
	"godiscordspeechbot/bot"
	"godiscordspeechbot/utils"
	"log"
	"time"
)

var NewsCategories = []string{"updates", "events", "release-notes", "general"}

var ForumAPIPath = "/v1/forums"
var NewsAPIPath = "/news/"

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

func GetNews(b *bot.Bot, category string) NewsResponse {
	respString := utils.CurlGet(b.GetLostArkURL(), NewsAPIPath+category)

	var news NewsResponse

	err := json.Unmarshal(respString, &news)

	if err != nil {
		log.Print(err)
	}

	return news
}

func GetForumUpdates(b *bot.Bot) ForumResponse {
	respString := utils.CurlGet(b.GetLostArkURL(), ForumAPIPath)

	var frm ForumResponse
	err := json.Unmarshal(respString, &frm)

	if err != nil {
		log.Print(err)
	}

	return frm
}

func GetNewsItems(hashmap map[uint32]bool, response NewsResponse, todayOnly bool) []DiscordMsg {
	var msgs []DiscordMsg

	for _, datum := range response.Data {
		hash := utils.Hash(datum.URL)

		// If it's not in the hashmap, it's new
		if _, ok := hashmap[hash]; !ok {
			hashmap[hash] = true

			if todayOnly {
				today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, 0, 0, 0, 0, time.UTC)
				forumItemDate, _ := time.Parse("January 02, 2006", datum.PublishDate)

				if forumItemDate.Before(today) {
					log.Println("Skipping Item...")
					continue
				}
			}

			msgs = append(msgs, DiscordMsg{
				"News",
				datum.Title,
				datum.Description,
				datum.URL,
			})
		}
	}

	return msgs
}

func GetForumItems(hashmap map[uint32]bool, response ForumResponse, todayOnly bool) []DiscordMsg {
	var msgs []DiscordMsg

	for _, datum := range response.Data {
		hash := utils.Hash(datum.URL)

		// If it's not in the hashmap, it's new
		if _, ok := hashmap[hash]; !ok {
			hashmap[hash] = true

			if todayOnly {
				today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, 0, 0, 0, 0, time.UTC)
				forumItemDate := datum.CreatedAt

				if forumItemDate.Before(today) {
					log.Println("Skipping Item...")
					continue
				}
			}

			lang := whatlanggo.Detect(datum.PostBody).Lang

			if whatlanggo.Eng != lang {
				log.Println("Skipping non-english forum item.")
				continue
			}

			msgs = append(msgs, DiscordMsg{
				"Forum",
				datum.Title,
				datum.PostBody,
				datum.URL,
			})
		}
	}
	return msgs
}

func LostArkCog(b *bot.Bot, ctx *discordgo.MessageCreate, interval time.Duration) {
	newsHashMap := make(map[uint32]bool)
	forumHashap := make(map[uint32]bool)

	for {
		for range time.Tick(interval) {
			for _, category := range NewsCategories {
				news := GetNewsItems(newsHashMap, GetNews(b, category), b.TodayOnly())

				if len(news) > 0 {
					log.Println("News found!")

					for i := 0; i < len(news); i++ {
						newsItem := news[i]
						b.SendMsgChannel(ctx.ChannelID, newsItem.URL)
					}
				}
			}

			forumUpdates := GetForumItems(forumHashap, GetForumUpdates(b), b.TodayOnly())

			if len(forumUpdates) > 0 {
				log.Println("Forum updates found!")
				for i := 0; i < len(forumUpdates); i++ {
					forumUpdate := forumUpdates[i]
					b.SendMsgChannel(ctx.ChannelID, forumUpdate.URL)
				}
			}
		}
	}
}
