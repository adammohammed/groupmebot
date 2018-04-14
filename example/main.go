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

	loggerlist := make([]groupmebot.Logger, 2)
	loggerlist[0] = groupmebot.CSVLogger{LogFile: "test.csv"}
	loggerlist[1] = groupmebot.StdOutLogger{}

	compLogger := groupmebot.CompositeLogger{loggerlist}
	cfg, err := groupmebot.NewBotConfigFromJson("mybot_cfg.json")

	if err != nil {
		log.Fatal("Could not update bot structure")
	}

	bot := groupmebot.GroupMeBot{BotConfig: cfg, Logger: compLogger}

	log.Printf("%s : %s ", bot.LogFile, bot.ID)
	// Make a list of functions
	bot.AddHook("Hi!$", hello)
	bot.AddHook("Hello!$", hello2)

	// Create Server to listen for incoming POST from GroupMe
	log.Printf("Listening on %v...\n", bot.Server)
	http.HandleFunc("/", bot.Handler())
	log.Fatal(http.ListenAndServe(bot.Server, nil))
}
