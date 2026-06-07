//go:build js

package ui

import "syscall/js"

func OpenURL(url string) error {
	js.Global().Get("window").Call("open", url, "_blank", "noopener,noreferrer")
	return nil
}
