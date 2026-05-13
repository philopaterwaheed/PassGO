//go:build js && wasm

package storage

import (
	"encoding/json"
	"errors"
	"syscall/js"
)

const localStorageKey = "passgo.session"

type jsSessionStore struct{}

func newSessionStore() SessionStore {
	return &jsSessionStore{}
}

func (s *jsSessionStore) Load() (Session, error) {
	ls, ok := localStorage()
	if !ok {
		return Session{}, nil
	}
	v := ls.Call("getItem", localStorageKey)
	if v.IsNull() || v.IsUndefined() {
		return Session{}, nil
	}
	str := v.String()
	if str == "" {
		return Session{}, nil
	}

	var sess Session
	if err := json.Unmarshal([]byte(str), &sess); err != nil {
		return Session{}, err
	}
	return sess, nil
}

func (s *jsSessionStore) Save(sess Session) error {
	ls, ok := localStorage()
	if !ok {
		return errors.New("localStorage unavailable")
	}
	b, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	ls.Call("setItem", localStorageKey, string(b))
	return nil
}

func (s *jsSessionStore) Clear() error {
	ls, ok := localStorage()
	if !ok {
		return nil
	}
	ls.Call("removeItem", localStorageKey)
	return nil
}

func localStorage() (js.Value, bool) {
	global := js.Global()
	storage := global.Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return js.Value{}, false
	}
	return storage, true
}
