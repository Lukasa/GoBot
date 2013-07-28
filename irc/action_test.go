package irc

import (
	"github.com/Lukasa/GoBot/struc"
	"testing"
)

type WriteBuffer struct {
	Buffer []byte
}

func (b *WriteBuffer) Write(p []byte) (int, error) {
	b.Buffer = make([]byte, len(p))
	copy(b.Buffer, p)
	return len(p), nil
}

func TestLogAction(t *testing.T) {
	prefixes := []string{
		"lukasa!cory@lukasa.co.uk",
		"lukasa",
		"test!~test@test.org",
	}

	messages := []string{
		"Hi there guuise!",
		"Another test message!",
		"More! Exclamation! Points!",
	}

	args := make([]string, 0)
	kwargs := make(map[string]string)

	buffer := new(WriteBuffer)
	logger := LogAction(buffer)

	for i, _ := range prefixes {
		msg := struc.NewIRCMessage()
		msg.Prefix = prefixes[i]
		msg.Trailing = messages[i]

		resp := logger(msg, args, kwargs)

		if resp != nil {
			t.Errorf("Response is not nil.")
		}

		if len(buffer.Buffer) == 0 {
			t.Errorf("Failed to write %v", messages[i])
		}
	}
}

// Test the PrintAction function works as expected.
func TestPrintAction(t *testing.T) {
	action := PrintAction("You're doing good work, ${1}!")
	kwargs := make(map[string]string)
	args := [][]string{
		[]string{"", "Fishcake"},
		[]string{"", "Lukasa"},
		[]string{"", "Cats", "Dogs"},
	}
	responses := []string{
		"You're doing good work, Fishcake!",
		"You're doing good work, Lukasa!",
		"You're doing good work, Cats!",
	}

	for i, arg := range args {
		msg := struc.NewIRCMessage()
		msg.Response = false
		msg.Command = struc.PRIVMSG
		msg.Arguments = []string{"#python-requests"}

		response := action(msg, arg, kwargs)

		if response.Response {
			t.Errorf("Response %v should not be a response", i)
		}

		if response.Command != struc.PRIVMSG {
			t.Errorf("Response %v has incorrect Command: expected %v, got %v", i, struc.PRIVMSG, response.Command)
		}

		if len(response.Arguments) != 1 {
			t.Errorf("Response %v has incorrect number of arguments: %v", i, len(response.Arguments))
		}

		if response.Arguments[0] != "#python-requests" {
			t.Errorf("Response %v has incorrect argument: %v", i, response.Arguments[0])
		}

		if response.Trailing != responses[i] {
			t.Errorf("Response %v has incorrect trailer: expected %v, got %v", i, responses[i], response.Trailing)
		}
	}
}

// Test that the Pong function works as expected.
func TestPong(t *testing.T) {
	ping := struc.NewIRCMessage()
	ping.Arguments = []string{"adams.freenode.net"}
	out := make(chan *struc.IRCMessage)

	go Pong(ping, out)

	pong := <-out
	if pong.Response {
		t.Errorf("Pong should not return responses.\n")
	}

	if pong.Command != struc.PONG {
		t.Errorf("Unexpected command type: expected %v, got %v\n", struc.PONG, pong.Command)
	}

	if (len(pong.Arguments) != 1) || (pong.Arguments[0] != ping.Arguments[0]) {
		t.Errorf("Unexpected argument array: expected %v, got %v\n", ping.Arguments, pong.Arguments)
	}
}
