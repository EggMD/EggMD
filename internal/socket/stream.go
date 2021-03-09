package socket

import (
	"sync"

	"github.com/pkg/errors"
)

var stream *Stream

// Stream 包含下面所有的文档会话。
type Stream struct {
	sync.RWMutex

	// key: 文档 UID
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

func (s *Stream) newDocument(uid string) (*DocSession, error) {
	s.Lock()
	defer s.Unlock()
	doc, err := NewDocSession(uid)
	if err != nil {
		return nil, err
	}

	s.documents[uid] = doc

	go doc.AutoSaveRoutine()

	return doc, nil
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
