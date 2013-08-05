package irc

import (
	"bytes"
	"github.com/Lukasa/GoBot/struc"
)

// DispatchMessages loops infinitely, accepting parsed versions of received IRC messages and dispatching them to botscripts.
// These botscripts return any messages they want to send on the 'out' channel.
// This can be viewed as the 'main' loop in GoBot.
func DispatchMessages(in, out chan *struc.IRCMessage, scripts []Botscript) error {
	// Before we begin looping, set the botscripts running. Each one has its own dedicated input channel.
	chans := beginScripts(out, scripts)

	// Pull messages off the input channel and dispatch them to each botscript.
	for {
		msg, ok := <-in
		if !ok {
			break
		}

		// Currently, we drop responses on the floor.
		if msg.Response {
			continue
		}

		// Quickly turn around PINGs.
		if msg.Command == struc.PING {
			go Pong(msg, out)
			continue
		}

		// If this isn't a PRIVMSG, drop it as well.
		if msg.Command != struc.PRIVMSG {
			continue
		}

		for _, channel := range chans {
			channel <- msg
		}
	}

	return nil
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
