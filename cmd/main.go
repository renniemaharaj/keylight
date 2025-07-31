package main

import (
	"syscall"

	"github.com/renniemaharaj/keylight/internal/hook"
	"github.com/renniemaharaj/keylight/internal/overlay"
	"github.com/renniemaharaj/keylight/internal/prompt"
)

func main() {
	// Hide console window on launch
	var _, _, _ = syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole").Call()

	go hook.Start()
	prompt.ShowRunningStatus()
	overlay.Overylay_WINAPI(hook.GetEventChannel())
	select {}
}
