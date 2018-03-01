package web

import (
	"sync"
	"time"
)

type MemSessionServer struct {
	pool       map[string]*Session
	mux        *sync.RWMutex
	duration   time.Duration
	gcInterval time.Duration
}

func (this *Server) UseMemSession(duration time.Duration, gcInterval time.Duration) {
	memSessionServer := &MemSessionServer{
		pool:       make(map[string]*Session),
		mux:        new(sync.RWMutex),
		duration:   duration,
		gcInterval: gcInterval,
	}
	memSessionServer.startGC()
	this.UseSession(memSessionServer)
}

func (this *MemSessionServer) startGC() {
	ticker := time.NewTicker(this.gcInterval)
	go func() {
		for t := range ticker.C {
			for token, session := range this.pool {
				if session.IsExpired(t) {
					this.Del(token)
				}
			}
		}
	}()
}

func (this *MemSessionServer) Del(token string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	delete(this.pool, token)
}

func (this *MemSessionServer) Get(token string) *Session {
	session, ok := this.pool[token]
	if ok {
		session.Refresh()
		return session
	}
	return this.add(token)
}

func (this *MemSessionServer) Set(token string, session *Session) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.pool[token] = session
}

func (this *MemSessionServer) add(token string) *Session {
	session := NewSession(token, this.duration)
	this.Set(token, session)
	return session
}
