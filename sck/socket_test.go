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
	Buffer    []byte
	SetBuffer chan []byte
}

func (c *MockConn) Read(b []byte) (int, error) {
	data := <-c.SetBuffer
	copy(b, data)
	return len(b), nil
}

func (c *MockConn) Write(b []byte) (int, error) {
	c.Buffer = make([]byte, len(b))

	// This will block until the calling thread is ready for the read to happen.
	<-c.SetBuffer

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
	conn.SetBuffer = make(chan []byte)

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
		conn.SetBuffer <- msg

		// Sleep to give the other goroutine some time to make the copy. This slows the tests,
		// but we just have to accept that.
		time.Sleep(200)

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

// TestReceiver tests the sck.Receiver() goroutine.
func TestReceiver(t *testing.T) {
	// Set up our variables.
	dataChan := make(chan []byte)
	conn := new(MockConn)
	conn.SetBuffer = make(chan []byte)

	// Prepare the messages we're going to receive.
	messages := [][]byte{
		[]byte("PASS secretpasswordhere"),
		[]byte("SERVICE dict * *.fr 0 0 :French Dictionary"),
		[]byte(":syrk!kalt@millennium.stealth.net QUIT :Gone to have lunch"),
	}

	// Start the goroutine.
	go Receiver(conn, dataChan)

	for _, msg := range messages {
		// Send the message in.
		conn.SetBuffer <- msg

		// Receieve the message back out.
		recvMsg := <-dataChan
		recvMsg = bytes.TrimRight(recvMsg, "\x00")

		// If they aren't the same, something horrible happened.
		if cmp := bytes.Compare(recvMsg, msg); cmp != 0 {
			t.Errorf(
				"Failed to read correctly: expected %v, got %v",
				msg,
				recvMsg)
		}
	}

	// Close the channel and confirm the tests don't panic.
	close(dataChan)
	time.Sleep(50 * time.Millisecond)
}
