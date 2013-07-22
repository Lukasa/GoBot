package irc

import (
	"bytes"
	"github.com/Lukasa/GoBot/struc"
)

// DispatchMessages loops infinitely, accepting parsed versions of received IRC messages and dispatching them to botscripts.
// These botscripts return any messages they want to send on the 'out' channel.
// This can be viewed as the 'main' loop in GoBot.
func DispatchMessages(in, out chan *struc.IRCMessage, scripts []Botscript) {
	// Before we begin looping, set the botscripts running. Each one has its own dedicated input channel.
	chans := beginScripts(out, scripts)

	// Pull messages off the input channel and dispatch them to each botscript.
	for {
		msg, ok := <-in
		if !ok {
			break
		}

		// Anything that isn't a PRIVMSG or a PING gets dropped.
		if (msg.Response) || ((msg.Command != struc.PRIVMSG) && (msg.Command != struc.PONG)) {
			continue
		}

		for _, channel := range chans {
			channel <- msg
		}
	}
}

// ParsingLoop provides a tight loop that pops values off the input channel and dispatches goroutines to parse them. This loop
// is very small to attempt to avoid bottlenecking.
func ParsingLoop(in chan []byte, out chan *struc.IRCMessage) {
	for {
		unparsed, ok := <-in
		if !ok {
			break
		}

		// We can get multiple messages in each packet.
		messages := bytes.Split(unparsed, []byte{'\r', '\n'})

		for _, message := range messages {
			if len(message) > 0 {
				go ParseIRCMessage(message, out)
			}
		}
	}
}

// UnParsingLoop provides a tight loop that pops values off the input channel and dispatches goroutines to unparse them. This
// loop is very small to avoid bottlenecking, as it will get quite a lot of traffic.
func UnParsingLoop(in chan *struc.IRCMessage, out chan []byte) {
	for {
		parsed, ok := <-in
		if !ok {
			break
		}

		go UnparseIRCMessage(parsed, out)
	}
}

// beginScripts starts all the botscripts executing, and returns a slice of channels that will send messages to those
// scripts.
func beginScripts(outchan chan *struc.IRCMessage, scripts []Botscript) []chan *struc.IRCMessage {
	chans := make([]chan *struc.IRCMessage, 0)

	for _, script := range scripts {
		newChan := make(chan *struc.IRCMessage)
		go script(newChan, outchan)
		chans = append(chans, newChan)
	}

	return chans
}
