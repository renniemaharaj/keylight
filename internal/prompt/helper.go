package prompt

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func openBrowser(url string) {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecuteW := shell32.NewProc("ShellExecuteW")

	// HWND = 0, Operation = "open", File = url
	shellExecuteW.Call(0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("open"))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(url))),
		0, 0, 1) // SW_SHOWNORMAL
}
