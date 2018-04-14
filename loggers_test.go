package groupmebot

import (
	"bufio"
	"fmt"
	"os"
)

func ExampleCSVWritten() {
	lg := CSVLogger{"test.csv"}
	msg := InboundMessage{Sender_id: "120123", Text: "This is a, check", Name: "Adam Code"}
	lg.LogMessage(msg)

	f, err := os.OpenFile(lg.LogFile, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("Couldn't open file\n")
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	text, _, err := rd.ReadLine()
	if err != nil {
		fmt.Printf("Couldn't open file\n")
	}

	os.Remove(lg.LogFile)
	fmt.Printf("Result: %s\n", text)

	// Output:
	// Result: 120123,"This is a, check",Adam Code
}
