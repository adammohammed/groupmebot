package groupmebot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type GroupMeBot struct {
	ID               string `json:"bot_id"`
	GroupID          string `json:"group_id"`
	Host             string `json:"host"`
	Port             string `json:"port"`
	LogFile          string `json:"logfile"`
	LogMethod        string `json:"logmethod"`
	Server           string
	TrackBotMessages bool `json:"trackbotmessages"`
	Hooks            map[string]func(InboundMessage) string
	Logger
}

type InboundMessage struct {
	Id           string                   `json:"id"`
	Avatar_url   string                   `json:"avatar_url"`
	Name         string                   `json:"name"`
	Sender_id    string                   `json:"sender_id"`
	Sender_type  string                   `json:"sender_type"`
	System       bool                     `json:"system"`
	Text         string                   `json:"text"`
	Source_guid  string                   `json:"source_guid"`
	Created_at   int                      `json:"created_at"`
	User_id      string                   `json:"user_id"`
	Group_id     string                   `json:"group_id"`
	Favorited_by []string                 `json:"favorited_by"`
	Attachments  []map[string]interface{} `json:"attachments"`
}

type OutboundMessage struct {
	ID   string `json:"bot_id"`
	Text string `json:"text"`
}

type Logger interface {
	LogMessage(msg InboundMessage)
}

/// NewBotFromJson (json cfg file name)
/// This reads a json file containing the keys
/// See the example bot_cfg.json
/// Returns err from ioutil if file can not be read
func NewBotFromJson(filename string) (*GroupMeBot, error) {
	file, err := ioutil.ReadFile(filename)

	var bot GroupMeBot
	if err != nil {
		log.Fatal("Error reading bot configuration json file")
		return nil, err
	}

	// Parse out information from file
	json.Unmarshal(file, &bot)

	bot.Server = bot.Host + ":" + bot.Port
	log.Printf("Creating bot at %s\nLogging at %s\n", bot.Server, bot.LogFile)
	bot.Hooks = make(map[string]func(InboundMessage) string)

	return &bot, err
}

func (b *GroupMeBot) SendMessage(outMessage string) (*http.Response, error) {
	msg := OutboundMessage{b.ID, outMessage}
	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	j_payload := string(payload)
	return http.Post("https://api.groupme.com/v3/bots/post", "application/json", strings.NewReader(j_payload))
}

func (b *GroupMeBot) AddHook(trigger string, response func(InboundMessage) string) {
	b.Hooks[trigger] = response
}

func (b *GroupMeBot) HandleMessage(msg InboundMessage) {
	resp := ""
	for trig, hook := range b.Hooks {
		matched, err := regexp.MatchString(trig, msg.Text)

		if matched {
			resp = hook(msg)
		} else if err != nil {
			log.Fatal("Error matching:", err)
		}

	}
	if len(resp) > 0 {
		log.Printf("Sending message: %v\n", resp)
		_, err := b.SendMessage(resp)
		if err != nil {
			log.Fatal("Error when sending.", err)
		}
	}

}

/*
 This is legitimate black magic, this is pretty cool, not usually able to do
 things like this in other languages. This is a function that takes
 a list of trigger functions and returns a function that can handle the Server
 Requests
*/
func (b *GroupMeBot) Handler() http.HandlerFunc {
	// Request Handler function

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			//log.Println("Bot recieving and handling message.")
			defer req.Body.Close()
			var msg InboundMessage
			err := json.NewDecoder(req.Body).Decode(&msg)
			if err != nil {
				log.Println("Couldn't parse the request body")
				msg.Sender_type = "bot"
			}
			if msg.Sender_type != "bot" || b.TrackBotMessages {
				b.LogMessage(msg)
				// Find hook by running through hooklist
				b.HandleMessage(msg)
			}
		} else {
			//log.Println("Bot not responding to unknown message")
			//io.WriteString(w, "GOTEM")
		}
	}
}
