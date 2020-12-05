package socket

import (
	"sync"
	"time"

	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/ot/operation"
	"github.com/EggMD/EggMD/internal/ot/selection"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

var (
	ErrInvalidRevision = errors.New("ot/session: invalid revision")
)

// DocSession is a single websocket session.
type DocSession struct {
	sync.Mutex

	// Document data
	ShortID string
	Content string

	Clients []*Client // The connection clients

	Operations         []*operation.Operation
	LastModifiedUserID uint
	EventChan          chan ConnEvent
	Done               chan struct{}
}

func NewDocSession(shortID string, content string) *DocSession {
	return &DocSession{
		ShortID: shortID,
		Content: content,

		Mutex:     sync.Mutex{},
		Clients:   make([]*Client, 0),
		EventChan: make(chan ConnEvent),
		Done:      make(chan struct{}),
	}
}

func (d *DocSession) AutoSaveRoutine() {
	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-tick.C:
			d.Save()
		case <-d.Done:
			close(d.Done)
			log.Trace("Stop auto save routine: %v", d.ShortID)
			return
		}
	}
}

func (d *DocSession) appendClient(client *Client) {
	d.Lock()
	defer d.Unlock()

	d.Clients = append(d.Clients, client)
	d.SendClients()
}

func (d *DocSession) removeClient(client *Client) {
	d.Lock()
	defer d.Unlock()

	if len(d.Clients) == 0 {
		return
	}

	for k, c := range d.Clients {
		if c == client {
			d.Clients = append(d.Clients[:k], d.Clients[k+1:]...)
		}
	}

	// If the last client left, destroy the session and save the document.
	if len(d.Clients) == 0 {
		d.Done <- struct{}{}
		d.Save()
		stream.removeDocument(d.ShortID)
	}

	d.SendClients()
}

func (d *DocSession) BroadcastExcept(client *Client, msg *EventMessage) {
	for _, c := range d.Clients {
		if c != client {
			c.out <- msg
		}
	}
}

func (d *DocSession) Broadcast(msg *EventMessage) {
	for _, c := range d.Clients {
		c.out <- msg
	}
}

func (d *DocSession) SetSelection(client *Client, sel *selection.Selection) {
	if client != nil {
		client.Selection = *sel
	}
}

func (d *DocSession) AddOperation(revision int, op *operation.Operation) (*operation.Operation, error) {
	if revision < 0 || len(d.Operations) < revision {
		return nil, ErrInvalidRevision
	}
	// find concurrent operations client isn't yet aware of
	otherOps := d.Operations[revision:]

	// transform given operation against these operations
	for _, otherOp := range otherOps {
		op1, _, err := operation.Transform(op, otherOp)
		if err != nil {
			return nil, err
		}
		if op.Meta != nil {
			if m, ok := op.Meta.(*selection.Selection); ok {
				op1.Meta = m.Transform(otherOp)
			}
		}

		op = op1
	}

	// apply transformed op on the doc
	doc, err := op.Apply(d.Content)
	if err != nil {
		return nil, err
	}

	d.Lock()
	defer d.Unlock()
	d.Content = doc
	d.Operations = append(d.Operations, op)

	return op, nil
}

func (d *DocSession) Save() {
	log.Trace("Save document: %v", d.ShortID)
	opt := db.UpdateDocOptions{
		Content:            d.Content,
		LastModifiedUserID: d.LastModifiedUserID,
	}
	_ = db.Documents.UpdateByShortID(d.ShortID, opt)
}

func (d *DocSession) SendClients() {
	d.Broadcast(&EventMessage{
		Name: "clients",
		Data: d.Clients,
	})
}
