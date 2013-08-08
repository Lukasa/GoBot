package main

import (
	"flag"
	"fmt"
	"github.com/Lukasa/GoBot/irc"
	"github.com/Lukasa/GoBot/sck"
	"github.com/Lukasa/GoBot/struc"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
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
	sigs := make(chan os.Signal)

	// Parallelise.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse this string into an IRC server.
	server, err := struc.NewIRCServerFromHostnamePort(serverStr)
	if err != nil {
		fmt.Printf("Could not parse %v. Exiting.", serverStr)
		os.Exit(1)
	}

	conn, err := sck.Connect(server, sendChan, recvChan)
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

	// At this stage we have successfully started execution, so we can add our signal handling.
	signal.Notify(sigs, os.Interrupt, os.Kill)

	// Send a test registration just to prove we can.
	login(username, channel, sendChan)

	// Run forever, dispatching messages.
	go irc.DispatchMessages(parsingOut, unparsingIn, []irc.Botscript{logscript, printscript})

	// Block on receipt of a signal. We don't care what it is, just die.
	<-sigs

	// Close our channels. Begin with the channel for incoming messages in bytes, then follow the
	// loop
	close(recvChan)
	close(parsingOut)
	close(unparsingIn)
	close(sendChan)

	// At this stage everything should be stopped, so we can safely close the connection.
	err = (*conn).Close()
	if err != nil {
		fmt.Printf("An error occurred while closing the connection: %v\n", err)
		os.Exit(3)
	}

	// If anything is still running, we should be worried. Check the number of goroutines, and
	// if it's more than one (this one), dump them to stdout.
	if runtime.NumGoroutine() > 1 {
		fmt.Printf("Error: Outstanding goroutines!\n")
		stack := make([]byte, 100)
		written := runtime.Stack(stack, true)
		strStack := string(stack[0:written])
		fmt.Print(strStack)
		os.Exit(4)
	}

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
