package hook

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// The keyboard event channel for triggering overlay and integrating with key hook
var keyboardChan = make(chan types.KeyboardEvent, 100)

// Keyboard event channel get function
func GetEventChannel() chan types.KeyboardEvent {
	return keyboardChan
}

// Start function for installing hook
func Start() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// The local run function installs the key hook and subscribes to the keyboard event channel
func run() error {
	// Buffer size is depends on your need. The 100 is placeholder value.

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return err
	}

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("start capturing keyboard input")
	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("Received timeout signal")
			return nil
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return nil
		case k := <-keyboardChan:
			fmt.Printf("Received %v %v\n", k.Message, k.VKCode)
			continue
		}
	}
}
