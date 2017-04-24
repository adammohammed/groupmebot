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
	h := make([]func(IncomingMessage), 0, 4)

	// Add functions that will later be "hooked" into
	// as a callback when messages arrive from group chat
	h = append(h, hello1)
	h = append(h, hello2)

	// Create Server to listen for incoming POST from GroupMe
	log.Printf("Listening on %v...\n", bot.Server)
	http.HandleFunc("/", BotHandler)
	log.Fatal(http.ListenAndServe(bot.Server, nil))
}

// Dummy functions later the input will likely be an
// IncomingMessage struct instead of string
func hello1(msg IncomingMessage) {
	fmt.Println("Hello World")
}

func hello2(msg IncomingMessage) {
	fmt.Println("Hello,", msg.Name)
}

// Request Handler function
func BotHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		log.Println("Bot recieving and handling message.")
		defer req.Body.Close()
		var msg IncomingMessage
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			log.Fatal("Couldn't read all the body", err)
		}
	} else {
		log.Println("Bot not responding to unknown message")
		io.WriteString(w, "hello world.\n")
	}
}
