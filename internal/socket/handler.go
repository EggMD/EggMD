package socket

import (
	"time"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	log "unknwon.dev/clog/v2"
)

func Handler(ctx *context.Context, receiver <-chan *EventMessage, sender chan<- *EventMessage, done <-chan bool, disconnect chan<- int, errorChannel <-chan error) (int, string) {
	shortID := ctx.Params("shortid")
	doc, err := db.Documents.GetDocByShortID(shortID)
	if err != nil {
		return 404, "document not found"
	}

	stream := getStream()
	docSession, err := stream.getDocument(shortID)
	if err != nil {
		docSession = stream.newDocument(shortID, doc.Content)
	}

	client := &Client{
		ID:   ctx.User.ID,
		Name: ctx.User.Name,

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
			"document": docSession.content,
			"revision": len(docSession.Operations),
			"clients":  docSession.Clients,
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
