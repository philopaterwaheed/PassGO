//go:build android && !cgo

package ui

import "errors"

func OpenURL(url string) error {
	if NormalizeURL(url) == "" {
		return nil
	}
	return errors.New("open android url: cgo unavailable")
}
