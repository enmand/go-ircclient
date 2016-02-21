package irc

import (
	"crypto/tls"
	"fmt"
	"net"
)

const (
	// CONNECTED is the connected status of the client
	CONNECTED = "connected"

	// DISCONNECTED is the disconnected status of the client
	DISCONNECTED = "disconnected"
)

// Connect connects this client to the server given
func (i *Client) Connect(server string) error {
	i.Server = server

	return i.connect()
}

func (i *Client) connect() error {
	var err error
	i.events = make(chan *Event)

	if !i.TLS {
		i.conn, err = net.DialTimeout("tcp", i.Server, TIMEOUT)
	} else {
		i.conn, err = tls.DialWithDialer(&net.Dialer{
			Timeout: TIMEOUT,
		}, "tcp", i.Server, &tls.Config{
			InsecureSkipVerify: i.TLSVerify,
		})
	}
	if err != nil {
		return fmt.Errorf("Could not connect to server: %s", err)
	}

	i.events <- &Event{
		Command: CONNECTED,
	}

	return err

}

// Disconnect disconnects this client from the server it's connected to
func (i *Client) Disconnect() error {
	var err error

	err = i.Write(&Event{
		Command: IRC_QUIT,
	})

	i.disconnect()
	if err != nil {

		return err
	}

	close(i.events)

	return nil
}

func (i *Client) disconnect() {
	i.conn.Close()

	i.events <- &Event{
		Command: DISCONNECTED,
	}
}
