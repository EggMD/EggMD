package socket

import "github.com/EggMD/EggMD/internal/ot/selection"

type Client struct {
	// Client user data
	ID        uint
	Name      string
	Avatar    string
	Selection selection.Selection `json:"selection"`

	in         <-chan *EventMessage
	out        chan<- *EventMessage
	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
}

type ConnEvent struct {
	*EventMessage
	Client *Client
}
