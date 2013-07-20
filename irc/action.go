package irc

import (
	"fmt"
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
		logMsg := fmt.Sprintf("%v - %v: %v", now, sender, msg.Trailing)

		target.Write([]byte(logMsg))
		return nil
	}
}
