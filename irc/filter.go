package irc

import (
	"github.com/Lukasa/GoBot/struc"
	"regexp"
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

// RegexFilterFromRegex takes a regular expression and returns a filter function that applies the given regular expression
// to each incoming IRC message.
// For the moment the regular expressiont treats named groups just the same as normal groups.
func RegexFilterFromRegex(regex string) (Filter, error) {
	compiledRegex, err := regexp.Compile(regex)
	if err != nil {
		// We should log here.
		return nil, err
	}

	var out Filter = func(msg *struc.IRCMessage) (bool, []string, map[string]string) {
		kwargs := make(map[string]string)

		if msg == nil {
			return false, make([]string, 0), kwargs
		}

		args := compiledRegex.FindAllStringSubmatch(msg.Trailing, -1) // Match as many as possible
		if len(args) == 0 {
			return false, make([]string, 0), kwargs
		}

		// We throw away any subsequent matches, to make reasoning about arguments easier.
		return true, args[0], kwargs
	}

	return out, nil
}
