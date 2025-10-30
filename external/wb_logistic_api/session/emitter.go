package session

import (
	"sync"
)

type Event string

const (
	EventUpdateAccessToken  Event = "update_access_token"
	EventUpdateSessionToken Event = "update_merged_token"
	EventUpdateUserInfo     Event = "update_user_info"
	EventClear              Event = "clear"
)

type Callback func(session *Session)

type Emitter struct {
	mtx       sync.RWMutex
	callbacks map[Event][]Callback
}

func NewEmitter() *Emitter {
	return &Emitter{
		callbacks: make(map[Event][]Callback),
	}
}

func (e *Emitter) On(event Event, cb Callback) {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	e.callbacks[event] = append(e.callbacks[event], cb)
}

func (e *Emitter) Emit(event Event, session *Session) {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	for _, cb := range e.callbacks[event] {
		go cb(session)
	}
}
