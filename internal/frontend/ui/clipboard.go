package ui

import (
	"io"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
)

func CopyText(gtx layout.Context, text string) {
	if text == "" {
		return
	}
	gtx.Execute(clipboard.WriteCmd{
		Type: "application/text",
		Data: io.NopCloser(strings.NewReader(text)),
	})
}
