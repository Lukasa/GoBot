package main

import (
	"flag"
	"fmt"
	"github.com/Lukasa/GoBot/irc"
	"github.com/Lukasa/GoBot/sck"
	"github.com/Lukasa/GoBot/struc"
	"math/rand"
	"os"
	"strconv"
)

// main is the entry point for GoBot.
func main() {
	sendChan := make(chan []byte)
	recvChan := make(chan []byte)
	args := parseArgs()
	username := genUsername()
	serverStr := args[0]
	channel := args[1]

	// Parse this string into an IRC server.
	server, err := struc.NewIRCServerFromHostnamePort(serverStr)
	if err != nil {
		fmt.Printf("Could not parse %v. Exiting.", serverStr)
		os.Exit(1)
	}

	_, err = sck.Connect(server, sendChan, recvChan)
	if err != nil {
		fmt.Printf("Could not connect to %v:%v. Exiting.", server.IPAddr, server.Port)
		os.Exit(2)
	}

	// Prepare the botscripts. For this simple case we'll log everything, so add a YesFilter and a logger to stdout.
	writeAction := irc.LogAction(os.Stdout)
	regexFilter, _ := irc.RegexFilterFromRegex("!m (.*)")
	printAction := irc.PrintAction("You're doing truly excellent work, ${1}!")
	logscript := irc.BuildBotscript([]irc.Filter{irc.YesFilter}, []irc.Action{writeAction})
	printscript := irc.BuildBotscript([]irc.Filter{regexFilter}, []irc.Action{printAction})

	// We need a few extra channels. One from the parsing loop to the dispatch loop, one from the goroutines to the
	// unparsing loop.
	parsingOut := make(chan *struc.IRCMessage)
	unparsingIn := make(chan *struc.IRCMessage)

	// Set the loops going.
	go irc.ParsingLoop(recvChan, parsingOut)
	go irc.UnParsingLoop(unparsingIn, sendChan)

	// Send a test registration just to prove we can.
	login(username, channel, sendChan)

	// Run forever, dispatching messages.
	err = irc.DispatchMessages(parsingOut, unparsingIn, []irc.Botscript{logscript, printscript})

	return
}

// parseArgs parses the command line arguments and flags. Currently this is the world's most boring function, but
// I'll extend it as I go.
func parseArgs() []string {
	flag.Parse()
	args := flag.Args()
	return args
}

func genUsername() string {
	base := "GoBot-"
	uniqId := strconv.Itoa(int(rand.Int31()))
	return base + uniqId
}

// Send the messages needed to login
func login(username, channel string, out chan []byte) {
	nick := []byte(fmt.Sprintf("NICK %v\r\n", username))
	out <- nick

	user := []byte(fmt.Sprintf("USER %v 1 1 1 :%v\r\n", username, username))
	out <- user

	join := []byte(fmt.Sprintf("JOIN %v\r\n", channel))
	out <- join
}
