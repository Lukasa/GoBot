package irc

import (
	"github.com/Lukasa/GoBot/struc"
)

// A Botscript represents a series of filters and actions to be applied on receipt of an IRC message. These expect to be
// launched as goroutines that run for the life of the program. They can be stopped from running by closing the input channel.
// They will not close their output channel before they exit.
type Botscript func(in, out chan *struc.IRCMessage)
