package socket

import (
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/ot/operation"
	"github.com/EggMD/EggMD/internal/ot/selection"
)

const (
	JOIN       = "join"       // Broadcast to all the client connections that a new client attend.
	CLIENTS    = "clients"    // Online clients info
	REGISTERED = "registered" // Tell the sender it has been registered in this session.
	PERMISSION = "permission" // Document permission, it's sent when the document permission changed.
)

type EventMessage struct {
	Name string      `json:"e"`
	Data interface{} `json:"d,omitempty"`
}
type H map[string]interface{}

func respMessage(name string, data interface{}) *EventMessage {
	return &EventMessage{
		Name: name,
		Data: data,
	}
}

func handleEvent(doc *DocSession, client *Client, evt *EventMessage) {
	var canView, canEdit bool
	if client.UserID == doc.Document.OwnerID {
		canView = true
		canEdit = true
	} else {
		// Check permission
		canView, canEdit = doc.Document.HasPermission(client.UserID)
		if !canView {
			client.disconnect <- 0
			return
		}
	}

	switch evt.Name {
	case "join":
		client.out <- respMessage(REGISTERED, H{
			"client_id": client.ID,
			"user_id":   client.UserID,
		})
		doc.BroadcastExcept(client, respMessage(JOIN, H{
			"client_id": client.ID,
			"username":  client.Name,
		}))

		// Set permission
	case "permission":
		// Only owner can set the permission.
		if doc.Document.OwnerID != client.UserID {
			return
		}

		permission, ok := evt.Data.(float64)
		if !ok {
			return
		}
		doc.Document.Permission = uint(permission)
		err := db.Documents.SetPermission(doc.Document.UID, uint(permission))
		if err != nil {
			return
		}
		doc.Broadcast(respMessage(PERMISSION, permission))

		// Kick
		for _, client := range doc.Clients {
			if view, _ := doc.Document.HasPermission(client.UserID); !view && doc.Document.OwnerID != client.UserID {
				client.disconnect <- 0
			}
		}

	case "op":
		if !canEdit {
			return
		}

		// data: [revision, ops, selection?]
		data, ok := evt.Data.([]interface{})
		if !ok {
			break
		}
		if len(data) < 2 {
			break
		}
		// revision
		revf, ok := data[0].(float64)
		rev := int(revf)
		if !ok {
			break
		}
		// ops
		ops, ok := data[1].([]interface{})
		if !ok {
			break
		}
		top, err := operation.Unmarshal(ops)
		if err != nil {
			break
		}
		// selection (optional)
		if len(data) >= 3 {
			selm, ok := data[2].(map[string]interface{})
			if !ok {
				break
			}
			sel, err := selection.Unmarshal(selm)
			if err != nil {
				break
			}
			top.Meta = sel
		}

		top2, err := doc.AddOperation(rev, top)
		if err != nil {
			break
		}

		client.out <- respMessage("ok", nil)

		if sel, ok := top2.Meta.(*selection.Selection); ok {
			doc.SetSelection(client, sel)
			doc.BroadcastExcept(client, &EventMessage{"op", []interface{}{client.ID, top2.Marshal(), sel.Marshal()}})
		} else {
			doc.BroadcastExcept(client, &EventMessage{"op", []interface{}{client.ID, top2.Marshal()}})
		}

	case "sel":
		data, ok := evt.Data.(map[string]interface{})
		if !ok {
			break
		}
		sel, err := selection.Unmarshal(data)
		if err != nil {
			break
		}
		doc.SetSelection(client, sel)
		doc.BroadcastExcept(client, &EventMessage{"sel", []interface{}{client.ID, sel.Marshal()}})
	}

	doc.LastModifiedUserID = client.UserID
}
