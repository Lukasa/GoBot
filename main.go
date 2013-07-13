package main

import (
	"flag"
	"fmt"
	"github.com/Lukasa/GoBot/sck"
	"github.com/Lukasa/GoBot/struc"
)

// main is the entry point for GoBot.
func main() {
	sendChan := make(chan []byte)
	recvChan := make(chan []byte)
	args := parseArgs()
	serverStr := args[0]

	// Parse this string into an IRC server.
	server, err := struc.NewIRCServerFromHostnamePort(serverStr)
	if err != nil {
		fmt.Errorf("Could not parse %v. Exiting.", serverStr)
	}

	_, err = sck.Connect(server, sendChan, recvChan)
	if err != nil {
		fmt.Errorf("Could not connect to %v:%v. Exiting.", server.IPAddr, server.Port)
	}

	// Send a test registration just to prove we can.
	nick := []byte("NICK GoBot")
	sendChan <- nick
	resp := <-recvChan
	fmt.Println(string(resp))

	user := []byte("USER gobot 0 * :GoBot")
	sendChan <- user
	resp = <-recvChan
	fmt.Println(string(resp))

	return
}

// parseArgs parses the command line arguments and flags. Currently this is the world's most boring function, but
// I'll extend it as I go.
func parseArgs() []string {
	flag.Parse()
	args := flag.Args()
	return args
}
