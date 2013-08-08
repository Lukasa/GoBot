package irc

import (
	"bytes"
	"github.com/Lukasa/GoBot/struc"
	"sort"
	"strings"
)

// ParseIRCMessage takes a byte array of an IRC message and parses it into an IRCMessage structure.
// This function isn't necessarily cheap, so expects to be run as a goroutine.
// There are a ton of casual assumptions in here that make this function less-than-resilient. I'm going to call it
// "low-hanging fruit for future improvement" and move on with my life.
func ParseIRCMessage(msg []byte, out chan *struc.IRCMessage) {
	parsedMsg := struc.NewIRCMessage()
	response := false
	prefix := ""
	responseCode := ""
	command := -1
	arguments := make([]string, 0)
	newMsg := []byte{}
	trailer := ""

	// First, split the message on whitespace.
	components := bytes.Fields(msg)
	if len(components) == 0 {
		// All whitespace, nothing to parse.
		return
	}

	// Check whether the first component begins with a ':'. If it does, it's an IRC prefix.
	if components[0][0] == ':' {
		prefix = string(components[0][1:])
		components = components[1:]
	}

	if len(components) != 0 {
		// The new first component is either a command or a response code. Easiest way to tell is to work out
		// if the first character is actually a digit.
		if (components[0][0] <= '9') && (components[0][0] >= '0') {
			responseCode = string(components[0])
			response = true
		} else {
			// Find the index of the command string in the Commands array.
			commandStr := string(bytes.ToUpper(components[0]))
			command = sort.SearchStrings(struc.Commands, commandStr) // This isn't quite right if we get a command that we don't know.
		}

		components = components[1:]

		// It's possible the message ends there. If it doesn't, the remaining values are either arguments or the 'trailing'
		// section. The trailing section begins with a colon and the arguments shouldn't (I don't think), so use that as the
		// delimiter.
		for _, arg := range components {
			if arg[0] == ':' {
				break
			}

			arguments = append(arguments, string(arg))
		}

		// Anything left after this stage is part of the 'trailing' section. Join it up.
		components = components[len(arguments):]
		newMsg = bytes.Join(components, []byte{' '})

		// Strip the leading colon.
		if len(newMsg) > 0 {
			newMsg = newMsg[1:]
		}

		trailer = string(newMsg)
	}

	// Append everything to the IRCMessage.
	parsedMsg.Prefix = prefix
	parsedMsg.Response = response
	parsedMsg.Command = command
	parsedMsg.ResponseCode = responseCode
	parsedMsg.Arguments = arguments
	parsedMsg.Trailing = trailer

	// Send it on the channel and exit.
	out <- parsedMsg
	return
}

// UnparseIRCMessage builds an IRC message into a byte array suitable for sending on the wire. This method is not necessarily
// cheap, so it expects to be run as a goroutine.
func UnparseIRCMessage(msg *struc.IRCMessage, out chan []byte) {
	components := make([]string, 0, 5)

	if msg.Prefix != "" {
		components = append(components, ":"+msg.Prefix)
	}

	if msg.Response {
		components = append(components, msg.ResponseCode)
	} else {
		components = append(components, struc.Commands[msg.Command])
	}

	for _, arg := range msg.Arguments {
		components = append(components, arg)
	}

	if msg.Trailing != "" {
		components = append(components, ":"+msg.Trailing)
	}

	strMsg := strings.Join(components, " ") + "\r\n"
	out <- []byte(strMsg)

	return
}
