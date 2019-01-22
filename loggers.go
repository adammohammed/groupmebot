package groupmebot

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type CSVLogger struct {
	LogFile string
}

func (logger CSVLogger) LogMessage(msg InboundMessage) {
	id := fmt.Sprintf("%s", msg.Sender_id)
	txt := fmt.Sprintf("%s", msg.Text)
	name := fmt.Sprintf("%s", msg.Name)
	values := []string{id, txt, name}
	if  len(id) == 0 {
		return
	}

	f, err := os.OpenFile(logger.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal("Couldn't open file to log messages")
	}

	defer f.Close()
	fwriter := bufio.NewWriter(f)
	csvWriter := csv.NewWriter(fwriter)

	csvWriter.Write(values)
	csvWriter.Flush()
	fwriter.Flush()
}

type StdOutLogger struct {
}

func (logger StdOutLogger) LogMessage(msg InboundMessage) {
	id := fmt.Sprintf("%s", msg.Sender_id)
	txt := fmt.Sprintf("%s", msg.Text)
	name := fmt.Sprintf("%s", msg.Name)
	log.Printf("Received Message: %s [%s:%s]\n", txt, name, id)
}

type CompositeLogger struct {
	Loggers []Logger
}

func (suite CompositeLogger) LogMessage(msg InboundMessage) {
	for _, l := range suite.Loggers {
		l.LogMessage(msg)
	}
}
