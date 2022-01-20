package stadia

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"stadia2xbox/stadia/hid"
)

type Device struct {
	hid.Device
}

const stadiaVID, stadiaPID = 0x18D1, 0x9400

var Controllers = make(map[string]bool)

func Devices() []*hid.DeviceInfo {
	devices, _ := hid.Devices()
	a := []*hid.DeviceInfo{}
	for _, device := range devices {
		if device.VendorID == stadiaVID && device.ProductID == stadiaPID && !Controllers[device.Path] {
			a = append(a, device)
		}
	}
	return a
}

func Open(device *hid.DeviceInfo) (*Device, error) {
	reEnable(device.Path)
	d, err := device.Open()
	if err == nil {
		Controllers[device.Path] = true
	}
	return &Device{d}, err
}

func (d Device) Read() (*Controller, error) {
	buf, ok := <-d.ReadCh()
	if !ok || buf[0] != 3 || len(buf) < 10 {
		return nil, errors.New("Failed to read from device")
	}
	controller := parse(buf)
	return &controller, nil
}

func (c *Device) Vibrate(largeMotor, smallMotor byte) error {
	return c.Write([]byte{0x05, largeMotor, largeMotor, smallMotor, smallMotor})
}

func reEnable(path string) {
	var id = strings.Replace(path[strings.Index(path, `\\?\`)+4:strings.LastIndex(path, `#`)], `#`, `\`, -1)
	ps, _ := exec.LookPath("powershell.exe")
	cmd := exec.Command(ps, fmt.Sprintf(`Disable-PnpDevice -InstanceId "%s" -Confirm:0 ; Enable-PnpDevice -InstanceId "%s" -Confirm:0`, id, id))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}
