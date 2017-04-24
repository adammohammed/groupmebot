package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type GroupMeBot struct {
	ID      string `json:"bot_id"`
	GroupID string `json:"group_id"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Server  string
}

type IncomingMessage struct {
	Avatar_url  string `json:"avatar_url"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Sender_id   string `json:"sender_id"`
	Sender_type string `json:"sender_type"`
	System      bool   `json:"system"`
	Text        string `json:"text"`
	User_id     string `json:"user_id"`
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

	// Create server Mux
	return &bot, err
}

func main() {

	bot, err := NewBotFromJson("mybot_cfg.json")
	if err != nil {
		log.Fatal("Could not create bot structure")
	}
	fmt.Printf("The bot id is %v\nThe Group id is %v.\n", bot.ID, bot.GroupID)

	// Make a list of functions
	h := make([]func(IncomingMessage) (bool, string), 0, 4)
	h = append(h, hello)
	h = append(h, hello2)

	// Create Server to listen for incoming POST from GroupMe
	log.Printf("Listening on %v...\n", bot.Server)
	http.HandleFunc("/", BotHandler(h))
	log.Fatal(http.ListenAndServe(bot.Server, nil))
}

/*
 Test hook functions
 Each hook should match a certain string, and if it matches
 it should return a string of text
 Hooks will be traversed until match occurs
*/
func hello(msg IncomingMessage) (bool, string) {
	matched, err := regexp.MatchString("Hi!$", msg.Text)
	if err != nil {
		return matched, ""
	}
	resp := fmt.Sprintf("Hi, %v.", msg.Name)
	return matched, resp
}

func hello2(msg IncomingMessage) (bool, string) {
	matched, err := regexp.MatchString("Hello!$", msg.Text)
	if err != nil {
		return matched, ""
	}
	resp := fmt.Sprintf("Hello, %v.", msg.Name)
	return matched, resp
}

/*
 This is legitimate black magic, this is pretty cool, not usually able to do
 things like this in other languages. This is a function that takes
 a list of trigger functions and returns a function that can handle the Server
 Requests
*/
func BotHandler(hooks []func(IncomingMessage) (bool, string)) http.HandlerFunc {
	// Request Handler function
	numhooks := len(hooks)

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			log.Println("Bot recieving and handling message.")
			defer req.Body.Close()
			var msg IncomingMessage
			err := json.NewDecoder(req.Body).Decode(&msg)

			// Find hook by running through hooklist
			hook, resp := false, ""
			for i := 0; !hook && i < numhooks; i++ {
				hook, resp = hooks[i](msg)
			}

			if hook {
				log.Printf("Sending message: %v\n", resp)
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
