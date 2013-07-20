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
