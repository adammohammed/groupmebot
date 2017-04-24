package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type GroupMeBot struct {
	ID      string `json:"bot_id"`
	GroupID string `json:"group_id"`
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
	return &bot, err
}

func main() {

	bot, err := NewBotFromJson("bot_cfg.json")
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
}

// Dummy functions later the input will likely be an
// IncomingMessage struct instead of string
func hello1(data string) {
	fmt.Println("Hello World")
}

func hello2(data string) {
	fmt.Println("Hello,", data)
}
