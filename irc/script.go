package irc

import (
	"github.com/Lukasa/GoBot/struc"
)

// A Botscript represents a series of filters and actions to be applied on receipt of an IRC message. These expect to be
// launched as goroutines that run for the life of the program. They can be stopped from running by closing the input channel.
// They will not close their output channel before they exit.
type Botscript func(in, out chan *struc.IRCMessage)

func BuildBotscript(filters []Filter, actions []Action) Botscript {
	return func(in, out chan *struc.IRCMessage) {
		for {
			msg, ok := <-in

			// Closed channel, stop executing.
			if !ok {
				break
			}

			args := make([]string, 0)
			kwargs := make(map[string]string)
			pass := true

			// Apply the filters.
			for filter := range filters {
				pass, newargs, newkwargs := filter(msg)

				if !pass {
					break
				}

				append(args, newargs...)
				for k, v := range newkwargs {
					kwargs[k] = v
				}
			}

			// Don't do anything if any of the filters failed.
			if !pass {
				continue
			}

			// Apply the actions.
			for action := range actions {
				response := action(msg, args, kwargs)
				if response != nil {
					out <- response
				}
			}
		}
	}
}
