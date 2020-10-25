package xbox

/*
#include <stdint.h>
typedef struct
{
	uint16_t wButtons;
	uint8_t bLeftTrigger;
	uint8_t bRightTrigger;
	int16_t sThumbLX;
	int16_t sThumbLY;
	int16_t sThumbRX;
	int16_t sThumbRY;
} xusb_report;
*/
import "C"

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	dll                                  = windows.NewLazyDLL("ViGEmClient.dll")
	procAlloc                            = dll.NewProc("vigem_alloc")
	procFree                             = dll.NewProc("vigem_free")
	procConnect                          = dll.NewProc("vigem_connect")
	procDisconnect                       = dll.NewProc("vigem_disconnect")
	procTargetAdd                        = dll.NewProc("vigem_target_add")
	procTargetFree                       = dll.NewProc("vigem_target_free")
	procTargetRemove                     = dll.NewProc("vigem_target_remove")
	procTargetX360Alloc                  = dll.NewProc("vigem_target_x360_alloc")
	procTargetX360RegisterNotification   = dll.NewProc("vigem_target_x360_register_notification")
	procTargetX360UnregisterNotification = dll.NewProc("vigem_target_x360_unregister_notification")
	procTargetX360Update                 = dll.NewProc("vigem_target_x360_update")
)

func Open(onVibration func(vibration Vibration)) (*Emulator, error) {
	client, _, err := procAlloc.Call()
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	_, _, err = procConnect.Call(client)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	return &Emulator{client, onVibration}, nil
}

func (e *Emulator) Connect() (*Controller, error) {
	target, _, err := procTargetX360Alloc.Call()
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	_, _, err = procTargetAdd.Call(e.handle, target)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	notif := func(client, target uintptr, largeMotor, smallMotor, ledNumber byte) uintptr {
		e.onVibration(Vibration{largeMotor, smallMotor})
		return 0
	}
	callback := windows.NewCallback(notif)

	_, _, err = procTargetX360RegisterNotification.Call(e.handle, target, callback)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}

	return &Controller{target, e}, nil
}

func (e *Emulator) Close() {
	procFree.Call(e.handle)
}

func (c *Controller) Close() {
	procTargetX360UnregisterNotification.Call(c.handle)
	procTargetRemove.Call(c.emulator.handle, c.handle)
	procTargetFree.Call(c.handle)
}

func (c *Controller) Send(report *Report) error {
	_, _, err := procTargetX360Update.Call(c.emulator.handle, c.handle, uintptr(unsafe.Pointer(&report.native)))
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}
	return nil
}

type Emulator struct {
	handle      uintptr
	onVibration func(vibration Vibration)
}

type Controller struct {
	handle   uintptr
	emulator *Emulator
}

type Report struct {
	native C.xusb_report
}

type Vibration struct {
	LargeMotor byte
	SmallMotor byte
}

func (r *Report) SetButton(check bool, button int) {
	if check {
		r.native.wButtons |= 1 << button
	}
}

func (r *Report) SetStick(dir bool, x, y int16) {
	if dir {
		r.native.sThumbRX = C.int16_t(x)
		r.native.sThumbRY = C.int16_t(y)
	} else {
		r.native.sThumbLX = C.int16_t(x)
		r.native.sThumbLY = C.int16_t(y)
	}
}

func (r *Report) SetTrigger(dir bool, value byte) {
	if dir {
		r.native.bRightTrigger = C.uint8_t(value)
	} else {
		r.native.bLeftTrigger = C.uint8_t(value)
	}
}

const (
	DPadUp      = 0
	DPadDown    = 1
	DPadLeft    = 2
	DPadRight   = 3
	ButtonStart = 4
	ButtonBack  = 5
	StickLeft   = 6
	StickRight  = 7
	BumperLeft  = 8
	BumperRight = 9
	ButtonGuide = 10
	ButtonA     = 12
	ButtonB     = 13
	ButtonX     = 14
	ButtonY     = 15
)
