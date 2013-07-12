package sck

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// Define a mock Connection for testing purposes. This will read to and write from a byte
// buffer.
type MockConn struct {
	Buffer []byte
}

// Currently does nothing, expand on it later.
func (c *MockConn) Read(b []byte) (int, error) {
	return 0, nil
}

func (c *MockConn) Write(b []byte) (int, error) {
	c.Buffer = make([]byte, len(b))
	copy(c.Buffer, b)
	return len(b), nil
}

func (c *MockConn) Close() error {
	return nil
}

func (c *MockConn) LocalAddr() net.Addr {
	return nil
}

func (c *MockConn) RemoteAddr() net.Addr {
	return nil
}

func (c *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// TestSender tests the sck.Sender() goroutine. This test appears to have unsynchronised
// access to a buffer but it actually doesn't, because the channel blocks the other goroutine.
func TestSender(t *testing.T) {
	// Set up our variables.
	dataChan := make(chan []byte)
	conn := new(MockConn)

	// Prepare the messages we're going to send.
	messages := [][]byte{
		[]byte("PASS secretpasswordhere"),
		[]byte("SERVICE dict * *.fr 0 0 :French Dictionary"),
		[]byte(":syrk!kalt@millennium.stealth.net QUIT :Gone to have lunch"),
	}

	// Start the goroutine.
	go Sender(conn, dataChan)

	for _, msg := range messages {
		dataChan <- msg

		// Sleep to give the other goroutine time to wake up. This slows the tests, but we just
		// have to accept that.
		time.Sleep(300 * time.Millisecond)

		if cmp := bytes.Compare(conn.Buffer, msg); cmp != 0 {
			t.Errorf(
				"Failed to write correctly: expected %v, got %v.",
				msg,
				conn.Buffer)
		}
	}

	// Close the channel.
	close(dataChan)
}
