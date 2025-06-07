package utils

import (
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	RequestUUID      uuid.UUID
	ExpieryTimeStamp time.Time
	UserIP           net.IP
	Index            int
}

var (
	sessions sync.Map
	heap     = &SessionHeap{}
	heapMu   sync.Mutex
)

// Data type for quick pruning of outdated requests
type SessionHeap []*UserSession

func (h SessionHeap) Len() int           { return len(h) }
func (h SessionHeap) Less(i, j int) bool { return h[i].ExpieryTimeStamp.Before(h[j].ExpieryTimeStamp) }
func (h SessionHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i]; h[i].Index = i; h[j].Index = j }

func (h *SessionHeap) Push(x any) {
	n := len(*h)
	item := x.(*UserSession)
	item.Index = n
	*h = append(*h, item)
}

func (h *SessionHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func (h SessionHeap) Peek() *UserSession {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}

func CreateSession(ip net.IP) uuid.UUID {
	id := uuid.New()
	expiresAt := time.Now().Add(5 * time.Minute)

	session := &UserSession{
		RequestUUID:      id,
		UserIP:           ip,
		ExpieryTimeStamp: expiresAt,
	}

	sessions.Store(id, session)

	heapMu.Lock()
	heap.Push(session)
	heapMu.Unlock()

	return id
}
