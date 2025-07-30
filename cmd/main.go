package main

import (
	"github.com/renniemaharaj/keylight/internal/hook"
	"github.com/renniemaharaj/keylight/internal/overlay"
)

func main() {
	go hook.Start()
	overlay.StartOverlay(hook.GetEventChannel())
}
