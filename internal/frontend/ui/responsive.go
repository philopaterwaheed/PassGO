package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// IsCompact returns true for phone-sized layouts.
func IsCompact(gtx layout.Context) bool {
	return gtx.Constraints.Max.X <= gtx.Dp(unit.Dp(600))
}
