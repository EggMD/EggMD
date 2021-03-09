package socket

import (
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/ot/operation"
	"github.com/EggMD/EggMD/internal/ot/selection"
)

const (
	JOIN       = "join"       // 新客户端加入的消息。
	CLIENTS    = "clients"    // 在线客户端信息。
	REGISTERED = "registered" // 客户端连接成功消息。
	PERMISSION = "permission" // 文档权限改变消息。
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
	userPermission := doc.Document.HasPermission(client.UserID)
	// 检查当前用户的文档权限
	if !userPermission.CanRead() {
		client.disconnect <- 0
		return
	}

	switch evt.Name {

	case "join": // 新客户端加入
		client.out <- respMessage(REGISTERED, H{
			"client_id": client.ID,
			"user_id":   client.UserID,
			"read_only": !userPermission.CanWrite(),
		})
		doc.BroadcastExcept(client, respMessage(JOIN, H{
			"client_id": client.ID,
			"username":  client.Name,
		}))

	case "permission": // 设置权限
		// 只有文档作者可以设置权限
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
			log.Error("Failed to set permission: %v", err)
			return
		}
		doc.Broadcast(respMessage(PERMISSION, permission))

		// 权限改变后，踢出无权限读取权限的用户。
		for _, client := range doc.Clients {
			if !doc.Document.HasPermission(client.UserID).CanRead() {
				client.disconnect <- 0
			}
		}

	case "op": // 操作
		if !userPermission.CanWrite() {
			return
		}

		doc.EditedAfterSave = true

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
