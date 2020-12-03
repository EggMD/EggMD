package socket

import (
	"github.com/EggMD/EggMD/internal/ot/operation"
	"github.com/EggMD/EggMD/internal/ot/selection"
)

type EventMessage struct {
	Name string      `json:"e"`
	Data interface{} `json:"d,omitempty"`
}

func handleEvent(doc *DocSession, client *Client, evt *EventMessage) {
	switch evt.Name {
	case "join":
		client.out <- &EventMessage{"registered", client.ID}
		doc.BroadcastExcept(client, &EventMessage{"join", map[string]interface{}{
			"client_id": client.ID,
			"username":  client.Name,
		}})
		
	case "op":
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

		client.out <- &EventMessage{"ok", nil}

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
}
