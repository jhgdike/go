package web

import (
	"errors"
	"time"
)

const (
	_SESSION_META_KEY        = "session"
	_SESSION_SERVER_META_KEY = "session_server"
)

type Session struct {
	Token    string                 `json:"token"`
	meta     map[string]interface{} `json:"meta"`
	ExpireAt time.Time              `json:"expire_at"`
	Duration time.Duration          `json:"duration"`
}

type SessionServer interface {
	Get(string) *Session
	Set(string, *Session)
	Del(string)
}

func (this *Server) UseSession(sessionServer SessionServer) {
	this.sessionServer = sessionServer
	this.Use(func(c *Context) {
		c.metaInternal[_SESSION_SERVER_META_KEY] = this.sessionServer
		c.Next()
		if session, ok := c.metaInternal[_SESSION_META_KEY]; ok {
			s := session.(*Session)
			this.sessionServer.Set(s.Token, s)
		}
	})
}

// 支持自己实现的session方法
func (this *Context) NewSession(token string, duration time.Duration) {
	session := NewSession(token, duration)
	this.metaInternal[_SESSION_META_KEY] = session
}

func (this *Context) StartSession(token string) error {
	sessionServer, ok := this.metaInternal[_SESSION_SERVER_META_KEY]
	if !ok {
		return errors.New("this server dose not use any session server")
	}
	session := sessionServer.(SessionServer).Get(token)
	this.metaInternal[_SESSION_META_KEY] = session
	return nil
}

func (this *Context) GetSession(key string) (interface{}, error) {
	session, ok := this.metaInternal[_SESSION_META_KEY]
	if !ok {
		return nil, errors.New("you should start session before get")
	}
	return session.(*Session).Get(key), nil
}

func (this *Context) SetSession(key string, value interface{}) error {
	session, ok := this.metaInternal[_SESSION_META_KEY]
	if !ok {
		return errors.New("you should start session before set")
	}
	session.(*Session).Set(key, value)
	return nil
}

func (this *Context) DestroySession(token string) error {
	sessionServer, ok := this.metaInternal[_SESSION_SERVER_META_KEY]
	if !ok {
		return errors.New("this server dose not use any session server")
	}
	sessionServer.(SessionServer).Del(token)
	delete(this.metaInternal, _SESSION_META_KEY)
	return nil
}

func NewSession(token string, duration time.Duration) *Session {
	return &Session{
		Token:    token,
		meta:     make(map[string]interface{}),
		ExpireAt: time.Now().Add(duration),
		Duration: duration,
	}
}

func (this *Session) Get(key string) interface{} {
	if value, ok := this.meta[key]; ok {
		return value
	}
	return nil
}

func (this *Session) Set(key string, value interface{}) {
	this.meta[key] = value
}

func (this *Session) Refresh() {
	this.ExpireAt = time.Now().Add(this.Duration)
}

func (this *Session) IsExpired(t time.Time) bool {
	return t.After(this.ExpireAt)
}
