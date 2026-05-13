//go:build !js

package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type nativeSessionStore struct {
	path string
}

func newSessionStore() SessionStore {
	path := defaultSessionPath()
	return &nativeSessionStore{path: path}
}

func defaultSessionPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		configDir = "."
	}
	return filepath.Join(configDir, "passgo", "session.json")
}

func (s *nativeSessionStore) Load() (Session, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Session{}, nil
		}
		return Session{}, err
	}

	var sess Session
	if err := json.Unmarshal(b, &sess); err != nil {
		return Session{}, err
	}
	return sess, nil
}

func (s *nativeSessionStore) Save(sess Session) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}

	b, err := json.Marshal(sess)
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, b, 0o600)
}

func (s *nativeSessionStore) Clear() error {
	err := os.Remove(s.path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
