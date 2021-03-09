package socket

import "github.com/EggMD/EggMD/internal/ot/selection"

// Client 为一个 WebSocket 客户端连接。
type Client struct {
	// Client user data
	ID        string              `json:"id"`
	UserID    uint                `json:"user_id"`
	Name      string              `json:"name"`
	Avatar    string              `json:"avatar"`
	Selection selection.Selection `json:"selection"`

	in         <-chan *EventMessage
	out        chan<- *EventMessage
	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
}
