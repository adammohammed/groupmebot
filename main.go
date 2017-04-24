package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type GroupMeBot struct {
	ID      string `json:"bot_id"`
	GroupID string `json:"group_id"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Server  string
	Hooks   []func(InboundMessage) (bool, string)
}

type InboundMessage struct {
	Avatar_url  string `json:"avatar_url"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Sender_id   string `json:"sender_id"`
	Sender_type string `json:"sender_type"`
	System      bool   `json:"system"`
	Text        string `json:"text"`
	User_id     string `json:"user_id"`
}

type OutboundMessage struct {
	ID   string `json:"bot_id"`
	Text string `json:"text"`
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

	bot.Hooks = make([]func(InboundMessage) (bool, string), 0, 10)

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

/*
 This is legitimate black magic, this is pretty cool, not usually able to do
 things like this in other languages. This is a function that takes
 a list of trigger functions and returns a function that can handle the Server
 Requests
*/
func (b *GroupMeBot) Handler() http.HandlerFunc {
	// Request Handler function
	numhooks := len(b.Hooks)

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			log.Println("Bot recieving and handling message.")
			defer req.Body.Close()
			var msg InboundMessage
			err := json.NewDecoder(req.Body).Decode(&msg)

			// Find hook by running through hooklist
			hook, resp := false, ""
			for i := 0; !hook && i < numhooks; i++ {
				hook, resp = b.Hooks[i](msg)
			}

			if hook {
				log.Printf("Sending message: %v\n", resp)
				_, err := b.SendMessage(resp)
				if err != nil {
					log.Fatal("Error when sending.", err)
				}
			}

			if err != nil {
				log.Fatal("Couldn't read all the body", err)
			}
		} else {
			log.Println("Bot not responding to unknown message")
			io.WriteString(w, "GOTEM")
		}
	}
}

/*
 Test hook functions
 Each hook should match a certain string, and if it matches
 it should return a string of text
 Hooks will be traversed until match occurs
*/
func hello(msg InboundMessage) (bool, string) {
	matched, err := regexp.MatchString("Hi!$", msg.Text)
	if err != nil {
		return matched, ""
	}
	resp := fmt.Sprintf("Hi, %v.", msg.Name)
	return matched, resp
}

func hello2(msg InboundMessage) (bool, string) {
	matched, err := regexp.MatchString("Hello!$", msg.Text)
	if err != nil {
		return matched, ""
	}
	resp := fmt.Sprintf("Hello, %v.", msg.Name)
	return matched, resp
}

func main() {

	bot, err := NewBotFromJson("mybot_cfg.json")
	if err != nil {
		log.Fatal("Could not create bot structure")
	}

	// Make a list of functions
	bot.Hooks = append(bot.Hooks, hello)
	bot.Hooks = append(bot.Hooks, hello2)

	// Create Server to listen for incoming POST from GroupMe
	log.Printf("Listening on %v...\n", bot.Server)
	http.HandleFunc("/", bot.Handler())
	log.Fatal(http.ListenAndServe(bot.Server, nil))
}
