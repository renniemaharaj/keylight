package hook

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
)

var (
	internalKeyboardChan = make(chan types.KeyboardEvent, 100)
	externalKeyBoardChan = make(chan types.KeyboardEvent, 100)

	internalMouseChan = make(chan types.MouseEvent, 100)

	ctrlDown       = false
	overlayEnabled = true
)

// External channel for overlay use
func GetEventChannel() chan types.KeyboardEvent {
	return externalKeyBoardChan
}

func Start() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := keyboard.Install(nil, internalKeyboardChan); err != nil {
		return err
	}
	defer keyboard.Uninstall()

	if err := mouse.Install(nil, internalMouseChan); err != nil {
		return err
	}
	defer mouse.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("Started keyboard hook")

	for {
		select {
		case <-signalChan:
			fmt.Println("Exiting: interrupt signal")
			return nil
		case k := <-internalKeyboardChan:
			keyboardHandler(&k)
		case m := <-internalMouseChan:
			mouseHandler(&m)
		}

	}
}
