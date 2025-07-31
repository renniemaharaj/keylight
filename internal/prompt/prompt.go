package prompt

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	supportLink = "https://www.paypal.com/paypalme/newrennie"
	statusTitle = "[Keylight v1] Running (background)"
	statusBody  = `
	Keylight is now running in the background.

  		CTRL + 1 → Toggle On
  		CTRL + 2 → Toggle Off
  		CTRL + 3 → Exit App

	Would you like to support future development?

	Click [Yes] to support ❤️
	Click [No] to continue using it for free.
`
)

func ShowRunningStatus() {

	// Prepare message
	title := windows.StringToUTF16Ptr(statusTitle)
	body := windows.StringToUTF16Ptr(statusBody)

	sendPrompt(title, body)
}

func sendPrompt(title, body *uint16) {
	user32 := windows.NewLazySystemDLL("user32.dll")
	messageBoxW := user32.NewProc("MessageBoxW")

	// 0x04 = MB_YESNO, 0x40 = MB_ICONQUESTION, 0x00000100 = MB_DEFBUTTON2
	flags := 0x00000004 | 0x00000040 | 0x00000100

	ret, _, _ := messageBoxW.Call(0,
		uintptr(unsafe.Pointer(body)),
		uintptr(unsafe.Pointer(title)),
		uintptr(flags))

	const IDYES = 6
	if ret == IDYES {
		openBrowser(supportLink)
	}
}
