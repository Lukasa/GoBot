# GoBot: A modular IRC bot for the masses.

It occurred to me that I'm a professional software developer and open-source
guy, but that I had never written an IRC bot. That's a criminal oversight, so
it seemed like time to start. While I'm at it, why not break out the new
hotness: [Go](http://golang.org/).

## How To Use GoBot

Invoking GoBot is really easy:

    $ gobot <IRC server> <channel name>

For example:

    $ gobot chat.freenode.net #python-requests

GoBot will set itself up, join the channel, and then...do nothing. If you want
a bot that does something a bit smarter, you can log all the output from the
channel:

    $ gobot -l logfile chat.freenode.net #python-requests

If you want to really get smarter, we need to start talking about botscripts.

## Botscripts

GoBot is all about modularity. You shouldn't need to know how to write Go in
order to get GoBot to do what you want. For that reason, GoBot knows how to
follow a set of instructions called a 'botscript'.

A very simple botscript looks like this:

    [filter]
    regex="!m .*"

    [action]
    print="You're doing great work, $UNAME!"

A botscript is made up of two parts: a filter section and then an action
section. Each of these sections is made up of several sub-parts. These
sub-parts are evaluated in order. The 'filter' section is applied to an IRC
message. If the filter matches, the 'action' section is performed in response.

Botscripts have a fairly limited set of functionality on their own, but GoBot
can go so far as to call arbitrary executable code if you set up the right
botscript. See the botscript documentation for more.

## Building GoBot

To build GoBot from source, make sure you've
[installed the Go toolchain](http://golang.org/doc/install). Then, simply run

    $ go install github.com/Lukasa/GoBot

Go will magically handle the rest.

## Documentation

Still to come.

## Contributing

As always, I welcome contributors! Please follow these steps:

1. Check for any open _or closed_ issues discussing your particular feature or
   bug.
2. Fork the repository.
3. Write a test that reproduces the bug.
4. Write the fix.
5. Send your Pull Request! If I don't get around to at least commenting on it
   within a couple of days, badger me on Twitter or in IRC.
