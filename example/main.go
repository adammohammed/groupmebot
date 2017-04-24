package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/adammohammed/groupmebot"
)

/*
 Test hook functions
 Each hook should match a certain string, and if it matches
 it should return a string of text
 Hooks will be traversed until match occurs
*/
func hello(msg groupmebot.InboundMessage) (bool, string) {
	matched, err := regexp.MatchString("Hi!$", msg.Text)
	if err != nil {
		return matched, ""
	}
	resp := fmt.Sprintf("Hi, %v.", msg.Name)
	return matched, resp
}

func hello2(msg groupmebot.InboundMessage) (bool, string) {
	matched, err := regexp.MatchString("Hello!$", msg.Text)
	if err != nil {
		return matched, ""
	}
	resp := fmt.Sprintf("Hello, %v.", msg.Name)
	return matched, resp
}

func main() {

	bot, err := groupmebot.NewBotFromJson("mybot_cfg.json")
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
