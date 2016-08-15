// Package irc provides IRC client services in Golang
//
// About
//
// This package implements an simple IRC service, that can be used in Golang to
// build IRC clients, bots, or other tools.
//
// See also: https://tools.ietf.org/html/rfc2812
package irc

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// TIMEOUT is the connection timeout to the IRC server
const TIMEOUT = 1 * time.Minute

// IRC is the IRC client interface
type IRC interface {
	// Connect to an IRC server. Use the form address:port
	Connect(server string) error

	// Disconnect from an IRC server
	Disconnect() error

	// Read blocks while reading from the server
	Read() error

	Handle(f []Filter, h HandlerFunc)
	Write(ev *Event) error
}

// Client is the implementation of the IRC interface
type Client struct {
	// The client's nickname on the server
	Nick string

	// The client's Ident on the server
	Ident string

	// The client's hostname
	Host string

	// The client's masked hostname on the server (if masked)
	MaskedHost string

	// Server is the server name the Client is connecting to
	Server string

	// If this connection is a TLS connection
	TLS bool

	// Should this client verify the server's SSL certs
	TLSVerify bool

	// handlers for filtered events
	handlers []*Handler

	// Events broadcasted from the server
	events chan *Event

	// The network connection this client has to the server
	conn net.Conn
}

// NewClient returns a new IRC client
func NewClient(nick, ident string, tls, tlsverify bool) *Client {
	c := &Client{
		Nick:      nick,
		Ident:     ident,
		TLSVerify: tlsverify,
		TLS:       tls,
	}

	c.Handle(
		[]Filter{CommandFilter(IRC_PING)},
		func(ev *Event, c IRC) {
			ev.Command = IRC_PONG
			c.Write(ev)
		},
	)

	c.Handle(
		[]Filter{CommandFilter(CONNECTED)},
		func(ev *Event, r IRC) {
			c.authenticate(r)
		},
	)

	c.Handle(
		[]Filter{CommandFilter(IRC_ERR_NICKNAMEINUSE)},
		func(ev *Event, r IRC) {
			c.Nick = fixNick(c.Nick, r)
			writeNick(c.Nick, r)
		},
	)

	return c
}

func (i *Client) authenticate(c IRC) {
	writeNick(i.Nick, c)

	// RFC 2812 USER command
	c.Write(&Event{
		Command: IRC_USER,
		Parameters: []string{
			i.Ident,
			"0",
			"*",
			i.Nick,
		},
	})
}

func fixNick(nick string, c IRC) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	uniq := r.Int()

	newNick := fmt.Sprintf("%s_%d", nick, uniq)

	return newNick[:9] // minimum max length in 9
}

func writeNick(nick string, c IRC) {
	c.Write(&Event{
		Command: IRC_NICK,
		Parameters: []string{
			nick,
		},
	})

	timeout := make(chan bool)
	go func() {
		// Wait 400ms for an error
		time.Sleep(400 * time.Millisecond)
		timeout <- true
	}()
	<-timeout
}
