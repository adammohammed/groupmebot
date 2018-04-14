package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adammohammed/groupmebot"
)

/*
 Test hook functions
 Each hook should match a certain string, and if it matches
 it should return a string of text
 Hooks will be traversed until match occurs
*/
func hello(msg groupmebot.InboundMessage) string {
	resp := fmt.Sprintf("Hi, %v.", msg.Name)
	return resp
}

func hello2(msg groupmebot.InboundMessage) string {
	resp := fmt.Sprintf("Hello, %v.", msg.Name)
	return resp
}

func main() {

	// Create two channels for logging, one to csv, one to stdout
	lg := groupmebot.CSVLogger{LogFile: "test.csv"}
	stdout := groupmebot.StdOutLogger{}
	// Group the channels in a Composite Logger type
	combinedLogger := groupmebot.CompositeLogger{Loggers: []groupmebot.Logger{lg, stdout}}

	// Plug in the Loggers to the bot and configure with tokens etc.
	bot := groupmebot.GroupMeBot{Logger: combinedLogger}
	err := bot.ConfigureFromJson("mybot_cfg.json")

	if err != nil {
		log.Fatal("Could not update bot structure")
	}

	// Make a list of functions
	bot.AddHook("Hi!$", hello)
	bot.AddHook("Hello!$", hello2)

	// Create Server to listen for incoming POST from GroupMe
	log.Printf("Listening on %v...\n", bot.Server)
	http.HandleFunc("/", bot.Handler())
	log.Fatal(http.ListenAndServe(bot.Server, nil))
}
