package irc

import (
	"github.com/Lukasa/GoBot/struc"
)

// The Filter type is the type of the functions that act as filters. They accept IRCMessage structures and return
// whether or not the given message matches the filter.
type Filter func(*struc.IRCMessage) bool

// The YesFilter is the simplest kind of filter, allowing all non-nil messages to match.
func YesFilter(msg *struc.IRCMessage) bool {
	if msg != nil {
		return true
	} else {
		return false
	}
}
