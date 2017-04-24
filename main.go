package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type GroupMeBot struct {
	ID      string `json:"bot_id"`
	GroupID string `json:"group_id"`
	host    string `json:"host"`
	port    string `json:"port"`
	server  string
}

type IncomingMessage struct {
	avatar_url  string `json:"avatar_url"`
	id          string `json:"id"`
	name        string `json:"name"`
	sender_id   string `json:"sender_id"`
	sender_type string `json:"sender_type"`
	system      bool   `json:"system"`
	text        string `json:"text"`
	user_id     string `json:"user_id"`
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

	bot.host = "0.0.0.0"
	bot.port = ":8080"
	bot.server = bot.host + bot.port
	// Parse out information from file
	json.Unmarshal(file, &bot)

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
	h := make([]func(string), 0, 4)

	// Add functions that will later be "hooked" into
	// as a callback when messages arrive from group chat
	h = append(h, hello1)
	h = append(h, hello2)

	// Test running hooks with sample data
	// range returns 2 things, the first output is the index
	// and the second output is the value at that index
	for _, f := range h {
		f("Adam")
	}

	// Create Server to listen for incoming POST from GroupMe
	log.Printf("Listening on %v%v/...\n", bot.host, bot.port)
	http.HandleFunc("/", BotHandler)
	log.Fatal(http.ListenAndServe(bot.server, nil))
}

// Dummy functions later the input will likely be an
// IncomingMessage struct instead of string
func hello1(data string) {
	fmt.Println("Hello World")
}

func hello2(data string) {
	fmt.Println("Hello,", data)
}

// Request Handler function
func BotHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		log.Println("Bot received message")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal("Couldn't read all the body", err)
		}
		log.Println(string(body))
	} else {
		log.Println("Bot not responding to unknown message")
		io.WriteString(w, "hello world.\n")
	}
}
