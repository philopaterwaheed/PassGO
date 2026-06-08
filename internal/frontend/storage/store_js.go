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
	if ok {
		removePassGOStorageItems(ls)
	}

	ss := js.Global().Get("sessionStorage")
	if !ss.IsUndefined() && !ss.IsNull() {
		removePassGOStorageItems(ss)
	}

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

func removePassGOStorageItems(storage js.Value) {
	for i := storage.Get("length").Int() - 1; i >= 0; i-- {
		key := storage.Call("key", i)
		if key.IsNull() || key.IsUndefined() {
			continue
		}
		keyStr := key.String()
		if keyStr == localStorageKey || len(keyStr) >= len("passgo.") && keyStr[:len("passgo.")] == "passgo." {
			storage.Call("removeItem", keyStr)
		}
	}
}
