package irc

import (
	"github.com/Lukasa/GoBot/struc"
)

// Actions are functions that respond to a given IRC message with a different IRC message.
type Action func(*struc.IRCMessage) *struc.IRCMessage

// NoAction defines the basic action to take in receipt of a message: none.
func NoAction(msg *struc.IRCMessage) *struc.IRCMessage {
	return nil
}
