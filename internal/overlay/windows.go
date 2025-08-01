//go:build windows

package overlay

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/moutend/go-hook/pkg/types"
	"golang.org/x/sys/windows"
)

// ==== DLL Imports ====
var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

// ==== Procedures from DLLs ====
var (
	createWindowExW            = user32.NewProc("CreateWindowExW")
	defWindowProcW             = user32.NewProc("DefWindowProcW")
	registerClassExW           = user32.NewProc("RegisterClassExW")
	showWindow                 = user32.NewProc("ShowWindow")
	destroyWindow              = user32.NewProc("DestroyWindow")
	dispatchMessageW           = user32.NewProc("DispatchMessageW")
	getMessageW                = user32.NewProc("GetMessageW")
	translateMessage           = user32.NewProc("TranslateMessage")
	getSystemMetrics           = user32.NewProc("GetSystemMetrics")
	setLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	getModuleHandleW           = kernel32.NewProc("GetModuleHandleW")
)

// ==== Window Constants ====
const (
	WS_POPUP         = 0x80000000
	WS_VISIBLE       = 0x10000000
	WS_EX_TOPMOST    = 0x00000008
	WS_EX_LAYERED    = 0x00080000
	WS_EX_TOOLWINDOW = 0x00000080
	SW_SHOW          = 5
	SW_HIDE          = 0
	WM_DESTROY       = 0x0002
	SM_CYSCREEN      = 1
	LWA_ALPHA        = 0x2
)

// ==== Globals ====
var (
	hInstance    windows.Handle
	hWnd         windows.Handle
	currentLevel = 0 // Default to 50%
	fadeOutTimer *time.Timer
)

var opacityLevels = map[int]byte{
	0: 128, // 50%
	1: 192, // 75%
	2: 255, // 100%
}

const (
	overlayHeight     = 150
	overlayClassName  = "KeyLightOverlay"
	overlayWindowName = ""
)

// ==== Data Structures ====
type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

type MSG struct {
	HWnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// ==== Public Entry ====
func Overylay_WINAPI(keyChan <-chan types.KeyboardEvent) {
	go runOverlay()
	for range keyChan {
		showOverlay()
	}
}

func HideOverlay() {
	if fadeOutTimer != nil {
		fadeOutTimer.Stop()
		hideOverlay()
	}
}

func CycleLevel() {
	currentLevel++
	if currentLevel > 2 {
		currentLevel = 0
	}
	setLevel(opacityLevels[currentLevel])
}

// ==== Internal Functions ====
func runOverlay() {
	handle, _, _ := getModuleHandleW.Call(0)
	hInstance = windows.Handle(handle)

	className, _ := syscall.UTF16PtrFromString(overlayClassName)
	wndProc := syscall.NewCallback(wndProc)

	wc := WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEX{})),
		Style:         0,
		LpfnWndProc:   wndProc,
		HInstance:     hInstance,
		HbrBackground: windows.Handle(6), // COLOR_WINDOW + 1
		LpszClassName: className,
	}
	registerClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	windowName, _ := syscall.UTF16PtrFromString(overlayWindowName)
	screenHeight, _, _ := getSystemMetrics.Call(SM_CYSCREEN)
	yPos := int(screenHeight) - overlayHeight

	hWndRaw, _, _ := createWindowExW.Call(
		WS_EX_TOPMOST|WS_EX_TOOLWINDOW|WS_EX_LAYERED,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		WS_POPUP|WS_VISIBLE,
		0,
		uintptr(yPos),
		1920,
		overlayHeight,
		0, 0,
		uintptr(hInstance),
		0,
	)
	hWnd = windows.Handle(hWndRaw)
	showWindow.Call(uintptr(hWnd), SW_HIDE)

	var msg MSG
	for {
		ret, _, _ := getMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		dispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func showOverlay() {
	showWindow.Call(uintptr(hWnd), SW_SHOW)
	setLevel(opacityLevels[currentLevel])
	if fadeOutTimer != nil {
		fadeOutTimer.Stop()
	}

	fadeOutTimer = time.AfterFunc(4*time.Second, func() {
		hideOverlay()
	})
}

func hideOverlay() {
	showWindow.Call(uintptr(hWnd), SW_HIDE)
}

func setLevel(alpha byte) {
	setLayeredWindowAttributes.Call(
		uintptr(hWnd),
		0,
		uintptr(alpha),
		LWA_ALPHA,
	)
}

func wndProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == WM_DESTROY {
		destroyWindow.Call(uintptr(hwnd))
		return 0
	}
	ret, _, _ := defWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret
}
