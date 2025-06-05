package pubsub

import (
	"comments-system/internal/models"
	"context"
	"sync"
)

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan *models.Comment]struct{}
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string]map[chan *models.Comment]struct{}),
	}
}

func (ps *PubSub) Subscribe(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	ch := make(chan *models.Comment, 10)

	ps.mu.Lock()
	if _, ok := ps.subscribers[postID]; !ok {
		ps.subscribers[postID] = make(map[chan *models.Comment]struct{})
	}
	ps.subscribers[postID][ch] = struct{}{}
	ps.mu.Unlock()

	go func() {
		<-ctx.Done()
		ps.mu.Lock()
		delete(ps.subscribers[postID], ch)
		close(ch)
		ps.mu.Unlock()
	}()

	return ch, nil
}

func (ps *PubSub) Publish(postID string, comment *models.Comment) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if subs, ok := ps.subscribers[postID]; ok {
		for ch := range subs {
			select {
			case ch <- comment:
			default:

			}
		}
	}
}
