package services

import (
	"container/heap"
	"net"
	"sync"
	"time"

	"github.com/MisterNorwood/SugarCube-Server/internal/database"
	. "github.com/MisterNorwood/SugarCube-Server/internal/utils"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
)

var (
	sessions sync.Map
	sHeap    = &SessionHeap{}
	heapMu   sync.Mutex
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
	sHeap.Push(session)
	heapMu.Unlock()

	return id
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

type SiteGetRequestResponse struct {
	RequestUUID   uuid.UUID     `json:"RequestID"`
	RequestedSite database.Site `json:"Site"`
}

func (SessionManager SessionManager) CreateResponseGetSite(ip net.IP, site database.Site) SiteGetRequestResponse {
	return SiteGetRequestResponse{
		RequestUUID:   CreateSession(ip),
		RequestedSite: site,
	}

}
