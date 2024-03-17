package main

// Original https://gist.github.com/NaniteFactory/0bd94e84bbe939cda7201374a0c261fd
import (
	"syscall"
	"unsafe"
)

const MB_OK = 0x00000000

// MessageBox of Win32 API.
func MessageBox(title ErrTitle, caption ErrMsg) int {
	hwnd := 0
	msgCaptionPtr, _ := syscall.UTF16PtrFromString(string(caption))
	msgTitlePtr, _ := syscall.UTF16PtrFromString(string(title))
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(msgCaptionPtr)),
		uintptr(unsafe.Pointer(msgTitlePtr)),
		uintptr(MB_OK))

	return int(ret)
}
