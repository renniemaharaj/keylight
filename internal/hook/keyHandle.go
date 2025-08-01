package hook

import (
	"fmt"
	"os"

	"github.com/moutend/go-hook/pkg/types"
	"github.com/renniemaharaj/keylight/internal/overlay"
)

// Key handling functions
func keyboardHandler(k *types.KeyboardEvent) {
	switch k.VKCode {
	case types.VK_LCONTROL, types.VK_RCONTROL:
		ctrlDown = (k.Message == types.WM_KEYDOWN)
	case types.VK_1:
		if ctrlDown && k.Message == types.WM_KEYDOWN {
			overlayEnabled = true
			fmt.Println("[Overlay] Activated")
		}
	case types.VK_2:
		if ctrlDown && k.Message == types.WM_KEYDOWN {
			overlayEnabled = false
			overlay.HideOverlay()
			fmt.Println("[Overlay] Deactivated")
		}
	case types.VK_3:
		if ctrlDown && k.Message == types.WM_KEYDOWN {
			fmt.Println("[Overlay] Exit requested")
			os.Exit(0)
		}
	case types.VK_4:
		if ctrlDown && k.Message == types.WM_KEYDOWN {
			overlay.CycleLevel()
		}
	}

	// Forward event if overlay is enabled
	if overlayEnabled {
		select {
		case externalKeyBoardChan <- *k:
		default:
			// Drop if externalChan is full
		}
	}
}
