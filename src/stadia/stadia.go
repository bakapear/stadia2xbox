package stadia

import (
	"errors"

	"./hid"
)

type Device struct {
	hid.Device
}

const stadiaVID, stadiaPID = 0x18D1, 0x9400

func Open() (*Device, error) {
	devices, _ := hid.Devices()
	for _, device := range devices {
		if device.VendorID == stadiaVID && device.ProductID == stadiaPID {
			device, err := device.Open()
			return &Device{device}, err
		}
	}
	return nil, errors.New("No stadia controller devices were found")
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
