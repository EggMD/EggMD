package socket

import (
	"time"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/tool"
	log "unknwon.dev/clog/v2"
)

func Handler(ctx *context.Context, receiver <-chan *EventMessage, sender chan<- *EventMessage, done <-chan bool, disconnect chan<- int, errorChannel <-chan error) (int, string) {
	uid := ctx.Params("uid")
	doc, err := db.Documents.GetDocByUID(uid)
	if err != nil {
		return 404, "document not found"
	}

	stream := getStream()
	docSession, err := stream.getDocument(uid)
	if err != nil {
		docSession = stream.newDocument(uid, doc.Content)
	}

	client := &Client{
		ID:     ctx.User.ID,
		Name:   ctx.User.Name,
		Avatar: tool.AvatarLink(ctx.User.AvatarEmail),

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
			"document":   docSession.Content,
			"revision":   len(docSession.Operations),
			"clients":    docSession.Clients,
			"permission": doc.Permission,
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
