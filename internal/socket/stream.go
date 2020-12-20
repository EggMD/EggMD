package socket

import (
	"sync"

	"github.com/pkg/errors"
)

var stream *Stream

// Stream top level
type Stream struct {
	sync.RWMutex

	// key: Document UID
	documents map[string]*DocSession
}

func getStream() *Stream {
	if stream == nil {
		stream = &Stream{
			RWMutex:   sync.RWMutex{},
			documents: make(map[string]*DocSession),
		}
	}
	return stream
}

func (s *Stream) newDocument(uid, content string) *DocSession {
	s.Lock()
	defer s.Unlock()
	doc := NewDocSession(uid, content)
	s.documents[uid] = doc

	go doc.AutoSaveRoutine()

	return doc
}

func (s *Stream) getDocument(uid string) (*DocSession, error) {
	s.Lock()
	defer s.Unlock()
	doc, ok := s.documents[uid]
	if !ok {
		return nil, errors.New("document session not found")
	}
	return doc, nil
}

func (s *Stream) removeDocument(uid string) {
	s.Lock()
	defer s.Unlock()
	delete(s.documents, uid)
}
