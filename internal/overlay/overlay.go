package overlay

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/moutend/go-hook/pkg/types"
)

// The timeout timer variable
var fadeOutTimer *time.Timer

// Rectangle for drawing rectangle
var rect *canvas.Rectangle

// Start function subscribes the overlay to the keyboard event channel, showing and hiding
func StartOverlay(keyChan <-chan types.KeyboardEvent) {
	a := app.New()
	w := a.NewWindow("KeyLight")
	w.SetPadded(false)
	w.Resize(fyne.NewSize(1920, 100))
	w.SetFixedSize(true)

	rect = canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(1920, 100))
	rect.Hide()

	w.SetContent(container.NewWithoutLayout(rect))
	w.Show()

	go func() {
		for range keyChan {
			show()
		}
	}()

	a.Run()
}

// Overlay show function, hides after timeout
func show() {
	rect.Show()
	canvas.Refresh(rect)

	if fadeOutTimer != nil {
		fadeOutTimer.Stop()
	}

	fadeOutTimer = time.AfterFunc(4*time.Second, func() {
		rect.Hide()
		canvas.Refresh(rect)
	})
}
