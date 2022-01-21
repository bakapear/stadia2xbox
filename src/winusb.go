package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	NULL                                = uintptr(0)
	WM_DEVICECHANGE                     = 537
	HWND_MESSAGE                        = ^uintptr(2)
	DEVICE_NOTIFY_ALL_INTERFACE_CLASSES = 4
	DBT_DEVTYP_DEVICEINTERFACE          = 5
	DBT_DEVICEARRIVAL                   = 0x8000
	DBT_DEVICEREMOVECOMPLETE            = 0x8004
)

var (
	user32                      = syscall.NewLazyDLL("user32.dll")
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	pDefWindowProc              = user32.NewProc("DefWindowProcW")
	pCreateWindowEx             = user32.NewProc("CreateWindowExW")
	pGetModuleHandle            = kernel32.NewProc("GetModuleHandleW")
	pRegisterClassEx            = user32.NewProc("RegisterClassExW")
	pGetMessage                 = user32.NewProc("GetMessageW")
	pDispatchMessage            = user32.NewProc("DispatchMessageW")
	pRegisterDeviceNotification = user32.NewProc("RegisterDeviceNotificationW")
)

var HID_DEVICE_CLASS = GUID{
	0x745a17a0,
	0x74d3,
	0x11d0,
	[8]byte{0xb6, 0xfe, 0x00, 0xa0, 0xc9, 0x0f, 0x57, 0xda},
}

var GUID_DEVINTERFACE_USB_DEVICE = GUID{
	0xa5dcbf10,
	0x6530,
	0x11d2,
	[8]byte{0x90, 0x1f, 0x00, 0xc0, 0x4f, 0xb9, 0x51, 0xed},
}

type MSG struct {
	hWnd    syscall.Handle
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
}

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type DevBroadcastDevinterface struct {
	dwSize       uint32
	dwDeviceType uint32
	dwReserved   uint32
	classGuid    GUID
	szName       uint16
}

type Wndclassex struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   syscall.Handle
	Icon       syscall.Handle
	Cursor     syscall.Handle
	Background syscall.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     syscall.Handle
}

func WndProc(hWnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_DEVICECHANGE:
		if wParam == uintptr(DBT_DEVICEARRIVAL) {
			DeviceChange()
		}
		return 0
	default:
		ret, _, _ := pDefWindowProc.Call(uintptr(hWnd), uintptr(msg), uintptr(wParam), uintptr(lParam))
		return ret
	}
}

func init() {
	className, _ := syscall.UTF16PtrFromString("stadia2xbox")

	cb := syscall.NewCallback(WndProc)
	mh, _, _ := pGetModuleHandle.Call(0)

	wc := Wndclassex{
		WndProc:   cb,
		Instance:  syscall.Handle(mh),
		ClassName: className,
	}

	wc.Size = uint32(unsafe.Sizeof(wc))
	a, _, err := pRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))

	if a == 0 {
		fmt.Printf("RegisterClassEx failed: %v", err)
		return
	}

	c := uintptr(unsafe.Pointer(className))

	ret, _, err := pCreateWindowEx.Call(NULL, c, c, NULL, NULL, NULL, NULL, NULL, HWND_MESSAGE, NULL, NULL, NULL)

	if ret == 0 {
		fmt.Printf("CreateWindowEx failed: %v", err)
		return
	}
	hWnd := syscall.Handle(ret)

	var notificationFilter DevBroadcastDevinterface
	notificationFilter.dwSize = uint32(unsafe.Sizeof(notificationFilter))
	notificationFilter.dwDeviceType = DBT_DEVTYP_DEVICEINTERFACE
	notificationFilter.dwReserved = 0
	notificationFilter.classGuid = HID_DEVICE_CLASS
	notificationFilter.szName = 0
	ret, _, err = pRegisterDeviceNotification.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&notificationFilter)), DEVICE_NOTIFY_ALL_INTERFACE_CLASSES)
	if ret == 0 {
		fmt.Printf("RegisterDeviceNotification failed: %v", err)
		return
	}

	var msg MSG
	go func() {
		for {
			if GlobalStop {
				return
			}
			ret, _, _ := pGetMessage.Call(uintptr(unsafe.Pointer(&msg)), NULL, NULL, NULL)
			if ret == 0 {
				break
			}
			pDispatchMessage.Call((uintptr(unsafe.Pointer(&msg))))
		}
	}()
}
