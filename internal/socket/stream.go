package socket

import (
	"sync"

	"github.com/pkg/errors"
)

var stream *Stream

// Stream top level
type Stream struct {
	sync.RWMutex

	// key: Document shortID
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

func (s *Stream) newDocument(shortID, content string) *DocSession {
	s.Lock()
	defer s.Unlock()
	doc := NewDocSession(shortID, content)
	s.documents[shortID] = doc
	return doc
}

func (s *Stream) getDocument(shortID string) (*DocSession, error) {
	s.Lock()
	defer s.Unlock()
	doc, ok := s.documents[shortID]
	if !ok {
		return nil, errors.New("document session not found")
	}
	return doc, nil
}
