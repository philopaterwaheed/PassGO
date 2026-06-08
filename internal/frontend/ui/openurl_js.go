//go:build js

package ui

import "syscall/js"

func OpenURL(url string) error {
	url = NormalizeURL(url)
	if url == "" {
		return nil
	}

	window := js.Global().Get("window")
	navigator := window.Get("navigator")
	isTouch := false
	if maxTouchPoints := navigator.Get("maxTouchPoints"); maxTouchPoints.Type() == js.TypeNumber {
		isTouch = maxTouchPoints.Int() > 0
	}
	if window.Get("matchMedia").Truthy() {
		isTouch = isTouch || window.Call("matchMedia", "(pointer: coarse)").Get("matches").Bool()
	}

	if isTouch {
		window.Get("location").Call("assign", url)
		return nil
	}

	opened := window.Call("open", url, "_blank")
	if opened.Truthy() {
		opened.Set("opener", js.Null())
	} else {
		window.Get("location").Call("assign", url)
	}
	return nil
}
