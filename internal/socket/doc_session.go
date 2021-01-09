package socket

import (
	"sync"
	"time"

	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/mdutil"
	"github.com/EggMD/EggMD/internal/ot/operation"
	"github.com/EggMD/EggMD/internal/ot/selection"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

var (
	ErrInvalidRevision = errors.New("ot/session: invalid revision")
)

// DocSession is a document collaborative session.
type DocSession struct {
	sync.Mutex

	Document *db.Document

	Clients            []*Client // The connection clients
	LastModifiedUserID uint

	Operations []*operation.Operation
	EventChan  chan ConnEvent
	Done       chan struct{}
}

// NewDocSession returns a new document collaborative session of the document `uid`.
func NewDocSession(uid string) (*DocSession, error) {
	document, err := db.Documents.GetDocByUID(uid)
	if err != nil {
		return nil, err
	}
	return &DocSession{
		Mutex:    sync.Mutex{},
		Document: document,

		Clients:   make([]*Client, 0),
		EventChan: make(chan ConnEvent),
		Done:      make(chan struct{}),
	}, nil
}

// AutoSaveRoutine save the document into database every 5 seconds.
func (d *DocSession) AutoSaveRoutine() {
	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-tick.C:
			d.Save()
		case <-d.Done:
			close(d.Done)
			log.Trace("Stop auto save routine: %v", d.Document.UID)
			return
		}
	}
}

// appendClient add a new client connection.
func (d *DocSession) appendClient(client *Client) {
	d.Lock()
	defer d.Unlock()

	d.Clients = append(d.Clients, client)
	d.BroadcastClientsInfo()
}

// removeClient remove the client connection.
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
		stream.removeDocument(d.Document.UID)
	}

	d.BroadcastClientsInfo()
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
	// Find concurrent operations client isn't yet aware of.
	otherOps := d.Operations[revision:]

	// Transform given operation against these operations.
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

	// Apply transformed op on the doc.
	doc, err := op.Apply(d.Document.Content)
	if err != nil {
		return nil, err
	}

	d.Lock()
	defer d.Unlock()
	d.Document.Content = doc
	d.Operations = append(d.Operations, op)

	return op, nil
}

func (d *DocSession) Save() {
	log.Trace("Save document: %v", d.Document.UID)
	opt := db.UpdateDocOptions{
		Title:              mdutil.ParseTitle(d.Document.Content),
		Content:            d.Document.Content,
		LastModifiedUserID: d.LastModifiedUserID,
	}
	_ = db.Documents.UpdateByUID(d.Document.UID, opt)
}

func (d *DocSession) BroadcastClientsInfo() {
	d.Broadcast(respMessage(CLIENTS, d.Clients))
}
