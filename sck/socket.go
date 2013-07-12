/*
Package socket defines the methods used to interact with sockets in GoBot.

These methods create and then use TCP connection objects to send and receive
IRC messages. They do no message decoding or encoding, instead focusing on
sending and recieving messages as fast as possible.
*/
package sck

import "net"

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
