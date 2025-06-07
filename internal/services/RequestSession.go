package services

import (
	"container/heap"
	"sync"
	"time"

	. "github.com/MisterNorwood/SugarCube-Server/internal/utils"
	"github.com/go-co-op/gocron"
)

type SessionManager struct {
	sessions sync.Map
	heap     *SessionHeap
	heapMu   sync.Mutex
}

func NewSessionManager() *SessionManager {
	h := &SessionHeap{}
	heap.Init(h)

	return &SessionManager{
		sessions: sync.Map{},
		heap:     h,
	}
}

func StartPruner(sm *SessionManager) {
	s := gocron.NewScheduler(time.UTC)
	s.Every(3).Second().Do(func() {
		now := time.Now()

		for {
			sm.heapMu.Lock()
			if sm.heap.Len() == 0 || (*sm.heap)[0].ExpieryTimeStamp.After(now) {
				sm.heapMu.Unlock()
				break
			}

			expired := heap.Pop(sm.heap).(*UserSession)
			sm.heapMu.Unlock()

			sm.sessions.Delete(expired.RequestUUID)
		}

	})

	s.StartAsync()
}
