package irc

import (
	"github.com/Lukasa/GoBot/struc"
)

// A Botscript represents a series of filters and actions to be applied on receipt of an IRC message. These expect to be
// launched as goroutines that run for the life of the program. They can be stopped from running by closing the input channel.
// They will not close their output channel before they exit.
type Botscript func(in, out chan *struc.IRCMessage)

// BuildBotscript creates a function that corresponds to a botscript that applies a variety of
// filters and actions to an incoming message.
func BuildBotscript(filters []Filter, actions []Action) Botscript {
	return func(in, out chan *struc.IRCMessage) {
		for {
			msg, ok := <-in

			// Closed channel, stop executing.
			if !ok {
				break
			}

			pass, args, kwargs := applyFilters(msg, filters)

			// Don't do anything if any of the filters failed.
			if !pass {
				continue
			}

			// Apply the actions.
			for _, action := range actions {
				response := action(msg, args, kwargs)
				if response != nil {
					out <- response
				}
			}
		}
	}
}

// applyFilters takes a map of filters and an IRC message and applies each filter in turn. Returns whether all the filters
// passed, and any arguments/keyword arguments they set.
func applyFilters(msg *struc.IRCMessage, filters []Filter) (bool, []string, map[string]string) {
	args := make([]string, 0)
	kwargs := make(map[string]string)
	pass := true

	for _, filter := range filters {
		var newargs []string
		var newkwargs map[string]string

		pass, newargs, newkwargs = filter(msg)

		if !pass {
			break
		}

		args = append(args, newargs...)
		for k, v := range newkwargs {
			kwargs[k] = v
		}
	}

	return pass, args, kwargs
}
