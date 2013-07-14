package struc

// Define the command types (good god there are a lot of these!)
const (
	ADMIN = iota
	AWAY
	CONNECT
	DIE
	ERROR
	INFO
	INVITE
	ISON
	JOIN
	KICK
	KILL
	LINKS
	LIST
	LUSERS
	MODE
	MOTD
	NAMES
	NICK
	NJOIN
	NOTICE
	OPER
	PART
	PASS
	PING
	PONG
	PRIVMSG
	QUIT
	REHASH
	RESTART
	SERVER
	SERVICE
	SERVLIST
	SQUERY
	SQUIRT
	SQUIT
	STATS
	SUMMON
	TIME
	TOPIC
	TRACE
	USER
	USERHOST
	USERS
	VERSION
	WALLOPS
	WHO
	WHOIS
	WHOWAS
)

// Commands defines the commands in string form. The use of iota above means the constants provide
// indices into this array.
var Commands = []string{
	"ADMIN",
	"AWAY",
	"CONNECT",
	"DIE",
	"ERROR",
	"INFO",
	"INVITE",
	"ISON",
	"JOIN",
	"KICK",
	"KILL",
	"LINKS",
	"LIST",
	"LUSERS",
	"MODE",
	"MOTD",
	"NAMES",
	"NICK",
	"NJOIN",
	"NOTICE",
	"OPER",
	"PART",
	"PASS",
	"PING",
	"PONG",
	"PRIVMSG",
	"QUIT",
	"REHASH",
	"RESTART",
	"SERVER",
	"SERVICE",
	"SERVLIST",
	"SQUERY",
	"SQUIRT",
	"SQUIT",
	"STATS",
	"SUMMON",
	"TIME",
	"TOPIC",
	"TRACE",
	"USER",
	"USERHOST",
	"USERS",
	"VERSION",
	"WALLOPS",
	"WHO",
	"WHOIS",
	"WHOWAS",
}

// IRCMessage structure represents a parsed IRC message. These are used internally rather than byte buffers.
type IRCMessage struct {
	Prefix       string   // With the colon removed.
	Response     bool     // Whether this is a response or a command.
	Command      int      // This or ResponseCode will be set, but not both.
	ResponseCode string   // This must be a string to ensure that 001 stays 001.
	Arguments    []string // Arbitrary number of arguments to any command.
	Trailing     string   // Anything that comes after the final colon in an IRC message.
}

// NewIRCMEssage builds and initialises a new IRC message.
func NewIRCMessage() *IRCMessage {
	msg := new(IRCMessage)
	msg.Arguments = make([]string, 0)
	return msg
}
