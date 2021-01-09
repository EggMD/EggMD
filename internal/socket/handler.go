package socket

import (
	"time"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/tool"
	"github.com/satori/go.uuid"
	log "unknwon.dev/clog/v2"
)

func Handler(ctx *context.Context, receiver <-chan *EventMessage, sender chan<- *EventMessage, done <-chan bool, disconnect chan<- int, errorChannel <-chan error) (int, string) {
	uid := ctx.Doc.UID

	stream := getStream()
	docSession, err := stream.getDocument(uid)
	if err != nil {
		docSession, err = stream.newDocument(uid)
		if err != nil {
			return 404, "document not found"
		}
	}

	var userID uint
	var name, avatar string

	if ctx.IsLogged {
		userID = ctx.User.ID
		name = ctx.User.Name
		avatar = tool.AvatarLink(ctx.User.Email)
	} else {
		userID = 0
		name = "Guest"
		avatar = tool.AvatarLink("")
	}

	client := &Client{
		ID:     uuid.NewV4().String(),
		UserID: userID,
		Name:   name,
		Avatar: avatar,

		in:         receiver,
		out:        sender,
		done:       done,
		err:        errorChannel,
		disconnect: disconnect,
	}
	docSession.appendClient(client)

	// Send document content, revision, connected clients.
	sender <- &EventMessage{
		"doc", map[string]interface{}{
			"document":   docSession.Document.Content,
			"revision":   len(docSession.Operations),
			"clients":    docSession.Clients,
			"owner_id":   docSession.Document.OwnerID,
			"permission": docSession.Document.Permission,
		},
	}

	ticker := time.After(30 * time.Minute)
	for {
		select {

		case evt := <-receiver:
			handleEvent(docSession, client, evt)
		case <-ticker:

		case <-done:
			docSession.removeClient(client)
			docSession.BroadcastExcept(client, &EventMessage{"quit", client.ID})

		case err := <-errorChannel:
			docSession.removeClient(client)
			docSession.BroadcastExcept(client, &EventMessage{"quit", client.ID})

			log.Error("connection error: %v", err)
			return 500, "an error occurred"
		}
	}
}
