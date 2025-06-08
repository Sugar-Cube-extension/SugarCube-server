package services

import (
	"container/heap"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/MisterNorwood/SugarCube-Server/internal/database"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
)

type UserSession struct {
	RequestUUID     uuid.UUID
	UserIP          net.IP
	ExpiryTimestamp time.Time
	index           int
}

type SessionHeap []*UserSession

func (h SessionHeap) Len() int           { return len(h) }
func (h SessionHeap) Less(i, j int) bool { return h[i].ExpiryTimestamp.Before(h[j].ExpiryTimestamp) }
func (h SessionHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *SessionHeap) Push(x any) {
	n := len(*h)
	session := x.(*UserSession)
	session.index = n
	*h = append(*h, session)
}

func (h *SessionHeap) Pop() any {
	old := *h
	n := len(old)
	session := old[n-1]
	session.index = -1 // for safety
	*h = old[0 : n-1]
	return session
}

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

func (sm *SessionManager) CreateSession(ip net.IP) uuid.UUID {
	id := uuid.New()
	expiry := time.Now().Add(5 * time.Minute)

	session := &UserSession{
		RequestUUID:     id,
		UserIP:          ip,
		ExpiryTimestamp: expiry,
	}

	sm.sessions.Store(id, session)

	sm.heapMu.Lock()
	heap.Push(sm.heap, session)
	sm.heapMu.Unlock()

	return id
}

func (sm *SessionManager) ValidateSession(id uuid.UUID) (bool, error) {
	val, ok := sm.sessions.Load(id)
	if !ok {
		return false, errors.New("session not found")
	}

	session := val.(*UserSession)
	if time.Now().After(session.ExpiryTimestamp) {
		sm.RemoveSession(id)
		return false, errors.New("session expired")
	}

	return true, nil
}

func (sm *SessionManager) RemoveSession(id uuid.UUID) {
	sm.sessions.Delete(id)
}

func (sm *SessionManager) StartPruner() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(3).Seconds().Do(func() {
		now := time.Now()

		for {
			sm.heapMu.Lock()
			if sm.heap.Len() == 0 || (*sm.heap)[0].ExpiryTimestamp.After(now) {
				sm.heapMu.Unlock()
				break
			}

			expired := heap.Pop(sm.heap).(*UserSession)
			sm.heapMu.Unlock()

			sm.sessions.Delete(expired.RequestUUID)
		}
	})
	scheduler.StartAsync()
}

func (sm *SessionManager) CreateResponseGetSite(ip net.IP, site database.Site) SiteGetRequestResponse {
	return SiteGetRequestResponse{
		RequestUUID:   sm.CreateSession(ip),
		RequestedSite: site,
	}
}

type SiteGetRequestResponse struct {
	RequestUUID   uuid.UUID     `json:"RequestID"`
	RequestedSite database.Site `json:"Site"`
}
