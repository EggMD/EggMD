package socket

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/mdutil"
	"github.com/EggMD/EggMD/internal/ot/operation"
	"github.com/EggMD/EggMD/internal/ot/selection"
)

var (
	ErrInvalidRevision = errors.New("ot/session: invalid revision")
)

// DocSession 是一个共享文档编辑会话。
type DocSession struct {
	sync.Mutex

	Document *db.Document

	Clients            []*Client // 当前连接的客户端
	LastModifiedUserID uint
	EditedAfterSave    bool

	Operations []*operation.Operation
	Done       chan struct{}
}

// NewDocSession 根据给定的 uid 返回一个对应共享文档的编辑会话。
func NewDocSession(uid string) (*DocSession, error) {
	document, err := db.Documents.GetDocByUID(uid)
	if err != nil {
		return nil, err
	}
	return &DocSession{
		Mutex:    sync.Mutex{},
		Document: document,

		Clients: make([]*Client, 0),
		Done:    make(chan struct{}),
	}, nil
}

// AutoSaveRoutine 每 5 秒保存一次文档内容到数据库中。
func (d *DocSession) AutoSaveRoutine() {
	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-tick.C:
			// 前 5 秒内有编辑操作
			if d.EditedAfterSave {
				d.Save()
			}
		case <-d.Done:
			close(d.Done)
			log.Trace("Stop auto save routine: %v", d.Document.UID)
			return
		}
	}
}

// appendClient 加入一个新的客户端连接。
func (d *DocSession) appendClient(client *Client) {
	d.Lock()
	defer d.Unlock()

	d.Clients = append(d.Clients, client)
	d.BroadcastClientsInfo()
}

// removeClient 移除一个新的客户端连接。
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

	// 如果最后一个客户端也离开了，则销毁该共享文档会话。
	if len(d.Clients) == 0 {
		d.Done <- struct{}{}
		d.Save()
		stream.removeDocument(d.Document.UID)
	}

	d.BroadcastClientsInfo()
}

// BroadcastExcept 向除了 client 以外的客户端广播消息。
func (d *DocSession) BroadcastExcept(client *Client, msg *EventMessage) {
	for _, c := range d.Clients {
		if c != client {
			c.out <- msg
		}
	}
}

// Broadcast 向所有客户单广播消息
func (d *DocSession) Broadcast(msg *EventMessage) {
	for _, c := range d.Clients {
		c.out <- msg
	}
}

// SetSelection 为 OT 算法文本选中操作。
func (d *DocSession) SetSelection(client *Client, sel *selection.Selection) {
	if client != nil {
		client.Selection = *sel
	}
}

// AddOperation 为 OT 算法新操作。
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

// Save 保存当前文档到数据库。
func (d *DocSession) Save() {
	log.Trace("Save document: %v", d.Document.UID)
	opt := db.UpdateDocOptions{
		Title:              mdutil.ParseTitle(d.Document.Content),
		Content:            d.Document.Content,
		LastModifiedUserID: d.LastModifiedUserID,
	}
	_ = db.Documents.UpdateByUID(d.Document.UID, opt)

	d.EditedAfterSave = false
}

// BroadcastClientsInfo 广播当前所有客户端信息。
func (d *DocSession) BroadcastClientsInfo() {
	d.Broadcast(respMessage(CLIENTS, d.Clients))
}
