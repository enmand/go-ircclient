package irc

// Responder
//
// Responder responds to the IRC server, by writing an IRC Event to the Client's
// connection

import (
	"bytes"
	"fmt"
)

// Write an event to the server, and return an error if it fails
func (i *Client) Write(ev *Event) error {
	var payload [][]byte

	payload = append(payload, []byte(ev.Command))
	for i, p := range ev.Parameters {
		if i == len(ev.Parameters)-1 && len(ev.Parameters) > 1 {
			p = fmt.Sprintf(":%s\r\n", p)
		}
		payload = append(payload, []byte(p))
	}

	payload = append(payload, []byte("\r\n"))
	full := bytes.Join(payload, []byte(" "))

	_, err := i.conn.Write(full)
	return err
}
