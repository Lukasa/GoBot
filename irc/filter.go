package irc

import (
	"github.com/Lukasa/GoBot/struc"
)

// The Filter type is the type of the functions that act as filters. They accept IRCMessage structures and return
// whether or not the given message matches the filter, and any arguments or keyword arguments they want passed to actions.
type Filter func(*struc.IRCMessage) (bool, []string, map[string]string)

// The YesFilter is the simplest kind of filter, allowing all non-nil messages to match.
func YesFilter(msg *struc.IRCMessage) (bool, []string, map[string]string) {
	args := make([]string, 0)
	kwargs := make(map[string]string)

	if msg != nil {
		return true, args, kwargs
	} else {
		return false, args, kwargs
	}
}
