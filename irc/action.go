package irc

import (
	"fmt"
	"github.com/Lukasa/GoBot/irc/util"
	"github.com/Lukasa/GoBot/struc"
	"io"
	"strings"
	"time"
)

// Actions are functions that respond to a given IRC message with a different IRC message.
type Action func(*struc.IRCMessage, []string, map[string]string) *struc.IRCMessage

// NoAction defines the basic action to take in receipt of a message: none.
func NoAction(msg *struc.IRCMessage, args []string, kwargs map[string]string) *struc.IRCMessage {
	return nil
}

// LogAction builds an action that writes any IRC message to a logging source (represented by the io.Writer).
func LogAction(target io.Writer) Action {
	return func(msg *struc.IRCMessage, args []string, kwargs map[string]string) *struc.IRCMessage {
		if msg == nil {
			return nil
		}

		now := time.Now().String()
		sender := strings.SplitN(msg.Prefix, "!", 2)[0]
		logMsg := fmt.Sprintf("%v - %v: %v\n", now, sender, msg.Trailing)

		target.Write([]byte(logMsg))
		return nil
	}
}

// PrintAction builds an action that replies to an IRC message with another message, based on a format
// string that's not unlike a shell string using variable substitution.
func PrintAction(format string) Action {
	// Get the format string and the argument indices.
	parsed, indices := util.BashFmtStringToGoFmtString(format)

	return func(msg *struc.IRCMessage, args []string, kwargs map[string]string) *struc.IRCMessage {
		if msg == nil {
			return nil
		}

		response := struc.NewIRCMessage()
		response.Response = false
		response.Command = struc.PRIVMSG
		response.Arguments = append(response.Arguments, msg.Arguments[0]) // First arg is the sender.

		// Build the response.
		fmtComponents := make([]interface{}, len(indices))
		for i, index := range indices {
			fmtComponents[i] = args[index]
		}

		response.Trailing = fmt.Sprintf(parsed, fmtComponents...)

		return response
	}
}

// Pong responds an IRC PING message with an IRC PONG message.
func Pong(msg *struc.IRCMessage, out chan *struc.IRCMessage) {
	outMsg := struc.NewIRCMessage()
	outMsg.Response = false
	outMsg.Command = struc.PONG

	// The PONG should basically return whatever was sent.
	for _, arg := range msg.Arguments {
		outMsg.Arguments = append(outMsg.Arguments, arg)
	}

	out <- outMsg
}
