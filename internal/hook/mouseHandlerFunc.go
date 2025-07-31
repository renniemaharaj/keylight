package hook

import (
	"github.com/moutend/go-hook/pkg/types"
	"github.com/renniemaharaj/keylight/internal/overlay"
)

// Mouse handling functions
func mouseHandler(_ *types.MouseEvent) {
	overlay.HideOverlay()
}
