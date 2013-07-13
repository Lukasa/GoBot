/*
Package socket defines the methods used to interact with sockets in GoBot.

These methods create and then use TCP connection objects to send and receive
IRC messages. They do no message decoding or encoding, instead focusing on
sending and recieving messages as fast as possible.
*/
package sck

import (
	"github.com/Lukasa/GoBot/struc"
	"net"
	"time"
)

// Sender loops indefinitely sending any data that is sent to it over the connection.
// Can be stopped by closing the channel.
// conn should be a net.TCPConn in real code, but has been left generic for testing purposes.
func Sender(conn net.Conn, data chan []byte) {
	for {
		msg, ok := <-data
		if !ok {
			break
		}

		_, err := conn.Write(msg)
		if err != nil {
			// Later we'll want to log this.
			break
		}
	}
}

// Receiver loops indefinitely reading data off a connection and passing it on the channel.
// Can be stopped by sending a close message on the 'cls' channel.
// conn should be a net.TCPConn in real code, but has been left generic for testing purposes.
func Receiver(conn net.Conn, data chan []byte) {
	// If the channel is forcefully closed this routine will panic (write on closed channel).
	// Allow that to happen, but don't kill the whole program.
	defer func() { recover() }()

	for {
		// Alloc a buffer into which we can read the data. According to RFC 2812:
		// "messages SHALL NOT exceed 512 characters in length".
		buf := make([]byte, 512)
		_, err := conn.Read(buf)
		if err != nil {
			// Later we'll want to log this.
			break
		}

		data <- buf
	}
}

// Connect sets up the connection to the IRC server and starts the goroutines that control sending
// and receiving data to/from the server. Returns the connection itself so that it can be closed
// at some later point.
func Connect(server struc.IRCServer, send, receive chan []byte) (*net.Conn, error) {
	addr := server.Name + ":" + string(server.Port)

	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return nil, err
	}

	go Sender(conn, send)
	go Receiver(conn, receive)

	return &conn, nil
}
