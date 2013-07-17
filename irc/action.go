package irc

import (
	"github.com/Lukasa/GoBot/struc"
)

// Actions are functions that respond to a given IRC message with a different IRC message.
type Action func(*struc.IRCMessage, []string, map[string]string) *struc.IRCMessage

// NoAction defines the basic action to take in receipt of a message: none.
func NoAction(msg *struc.IRCMessage, args []string, kwargs map[string]string) *struc.IRCMessage {
	return nil
}
