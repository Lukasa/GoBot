package irc

import (
	"bytes"
	"github.com/Lukasa/GoBot/struc"
	"testing"
)

func TestParseIRCMessage(t *testing.T) {
	// A set of pre-built IRC messages.
	messages := [][]byte{
		[]byte("USER GoBot 1 1 1 :GoBot"),
		[]byte("NICK GoBot"),
		[]byte(":barjavel.freenode.net NOTICE * :*** Checking Ident"),
		[]byte(":barjavel.freenode.net 002 GoBot :Your host is barjavel.freenode.net[78.40.125.4/6667], running version ircd-seven-1.1.3"),
		[]byte("PING calvino.freenode.net"),
	}

	prefixes := []string{
		"",
		"",
		"barjavel.freenode.net",
		"barjavel.freenode.net",
		"",
	}

	responses := []bool{
		false,
		false,
		false,
		true,
		false,
	}

	commands := []int{
		struc.USER,
		struc.NICK,
		struc.NOTICE,
		-1,
		struc.PING,
	}

	responseCodes := []string{
		"",
		"",
		"",
		"002",
		"",
	}

	arguments := [][]string{
		[]string{"GoBot", "1", "1", "1"},
		[]string{"GoBot"},
		[]string{"*"},
		[]string{"GoBot"},
		[]string{"calvino.freenode.net"},
	}

	trailers := []string{
		"GoBot",
		"",
		"*** Checking Ident",
		"Your host is barjavel.freenode.net[78.40.125.4/6667], running version ircd-seven-1.1.3",
		"",
	}

	recvChan := make(chan *struc.IRCMessage)

	// Spin up goroutines to parse the stuff.
	for _, msg := range messages {
		go ParseIRCMessage(msg, recvChan)
	}

	// Then get them out.
	for i := range messages {
		parsed := <-recvChan

		if prefixes[i] != parsed.Prefix {
			t.Errorf("Invalid prefix: expected %v, got %v", prefixes[i], parsed.Prefix)
		}
		if responses[i] != parsed.Response {
			t.Errorf("Invalid response value: expected %v, got %v", responses[i], parsed.Response)
		}
		if commands[i] != parsed.Command {
			t.Errorf("Invalid command value: expected %v, got %v", commands[i], parsed.Command)
		}
		if responseCodes[i] != parsed.ResponseCode {
			t.Errorf("Invalid response code: expected %v, got %v", responseCodes[i], parsed.ResponseCode)
		}
		for j, arg := range arguments[i] {
			if arg != parsed.Arguments[j] {
				t.Errorf("Invalid argument: expected %v, got %v", arg, parsed.Arguments[j])
			}
		}
		if trailers[i] != parsed.Trailing {
			t.Errorf("Invalid trailer: expected %v, got %v.", trailers[i], parsed.Trailing)
		}
	}
}

func TestUnparseIRCMessage(t *testing.T) {
	// A set of pre-built IRC messages.
	messages := [][]byte{
		[]byte("USER GoBot 1 1 1 :GoBot\r\n"),
		[]byte("NICK GoBot\r\n"),
		[]byte(":barjavel.freenode.net NOTICE * :*** Checking Ident\r\n"),
		[]byte(":barjavel.freenode.net 002 GoBot :Your host is barjavel.freenode.net[78.40.125.4/6667], running version ircd-seven-1.1.3\r\n"),
		[]byte("PING calvino.freenode.net\r\n"),
	}

	// The various components of each message.
	prefixes := []string{
		"",
		"",
		"barjavel.freenode.net",
		"barjavel.freenode.net",
		"",
	}

	responses := []bool{
		false,
		false,
		false,
		true,
		false,
	}

	commands := []int{
		struc.USER,
		struc.NICK,
		struc.NOTICE,
		-1,
		struc.PING,
	}

	responseCodes := []string{
		"",
		"",
		"",
		"002",
		"",
	}

	arguments := [][]string{
		[]string{"GoBot", "1", "1", "1"},
		[]string{"GoBot"},
		[]string{"*"},
		[]string{"GoBot"},
		[]string{"calvino.freenode.net"},
	}

	trailers := []string{
		"GoBot",
		"",
		"*** Checking Ident",
		"Your host is barjavel.freenode.net[78.40.125.4/6667], running version ircd-seven-1.1.3",
		"",
	}

	recvChan := make(chan []byte)

	// Loop over the built messages, build an IRCMessage structure corresponding to each, then unparse it.
	for i, msg := range messages {
		parsed := struc.NewIRCMessage()
		parsed.Prefix = prefixes[i]
		parsed.Response = responses[i]
		parsed.Command = commands[i]
		parsed.ResponseCode = responseCodes[i]
		parsed.Arguments = arguments[i]
		parsed.Trailing = trailers[i]

		go UnparseIRCMessage(parsed, recvChan)

		unparsed := <-recvChan

		if res := bytes.Compare(msg, unparsed); res != 0 {
			t.Errorf("Incorrectly unparsed: expected %v, got %v", msg, unparsed)
		}
	}
}
