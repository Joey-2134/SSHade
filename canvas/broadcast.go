package canvas

import "sync"

type Broadcaster struct {
	mu   sync.RWMutex
	subs []chan Pixel
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subs: make([]chan Pixel, 0),
	}
}

// Subscribe returns a receive-only channel for canvas updates and an unsubscribe function.
// Call the returned function when the session ends to avoid leaks.
func (b *Broadcaster) Subscribe() (<-chan Pixel, func()) {
	ch := make(chan Pixel, 16)
	b.mu.Lock()
	b.subs = append(b.subs, ch)
	b.mu.Unlock()
	unsub := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		for i, c := range b.subs {
			if c == ch {
				b.subs = append(b.subs[:i], b.subs[i+1:]...)
				close(ch)
				return
			}
		}
	}
	return ch, unsub
}

func (b *Broadcaster) Broadcast(p Pixel) {
	b.mu.RLock()
	subs := make([]chan Pixel, len(b.subs))
	copy(subs, b.subs)
	b.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- p:
		default:
			// subscriber slow or channel full, skip
		}
	}
}
