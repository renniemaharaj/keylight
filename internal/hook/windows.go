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

var (
	internalChan   = make(chan types.KeyboardEvent, 100)
	externalChan   = make(chan types.KeyboardEvent, 100)
	ctrlDown       = false
	overlayEnabled = true
)

// External channel for overlay use
func GetEventChannel() chan types.KeyboardEvent {
	return externalChan
}

func Start() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := keyboard.Install(nil, internalChan); err != nil {
		return err
	}
	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("Started keyboard hook")

	for {
		select {
		case <-time.After(10 * time.Minute):
			fmt.Println("Exiting: timed out")
			continue
		case <-signalChan:
			fmt.Println("Exiting: interrupt signal")
			return nil
		case k := <-internalChan:
			handleEvent(k)
		}
	}
}
