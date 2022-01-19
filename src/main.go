package main

import (
	"fmt"
	"stadia2xbox/stadia"
	"stadia2xbox/xbox"

	"github.com/getlantern/systray"
	"github.com/rodolfoag/gow32"
	"gopkg.in/toast.v1"
)

var stop = false

func main() {
	_, err := gow32.CreateMutex("stadia2xbox")
	if err != nil {
		msg("An instance of stadia2xbox is already running!")
	} else {
		systray.Run(ready, close)
	}
}

func ready() {
	icon, _ := Asset("data/stadia.ico")
	systray.SetIcon(icon)

	re := systray.AddMenuItem("Refresh (0 devices)", "Refresh devices")
	go repeat(re)

	quit := systray.AddMenuItem("Exit", "Exits the application")
	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()

	refresh(re)
}

func repeat(re *systray.MenuItem) {
	<-re.ClickedCh
	refresh(re)
	repeat(re)
}

func refresh(re *systray.MenuItem) {
	re.Disable()
	re.SetTitle("Refreshing...")

	connect(re)

	re.Enable()
	msg := fmt.Sprintf("Refresh (%d devices)", len(stadia.Controllers))
	re.SetTitle(msg)
}

func connect(re *systray.MenuItem) {
	device, err := stadia.Open()
	if err != nil {
		msg(err.Error())
		return
	}
	defer func() {
		device.Close()
		delete(stadia.Controllers, device.Info().Path)
		msg := fmt.Sprintf("Refresh (%d devices)", len(stadia.Controllers))
		re.SetTitle(msg)
	}()

	emu, err := xbox.Open(func(vibration xbox.Vibration) {
		device.Vibrate(vibration.LargeMotor, vibration.SmallMotor)
	})
	if err != nil {
		msg(err.Error())
		return
	}

	con, err := emu.Connect()
	if err != nil {
		msg(err.Error())
		return
	}
	defer emu.Close()
	defer con.Close()

	msg("Stadia Controller sucessfully connected and emulated as Xbox Controller")

	for {
		if stop {
			return
		}
		d, err := device.Read()
		if err != nil {
			msg(err.Error())
			return
		}
		report := xbox.Report{}

		report.SetButton(d.DPad.Up, xbox.DPadUp)
		report.SetButton(d.DPad.Down, xbox.DPadDown)
		report.SetButton(d.DPad.Left, xbox.DPadLeft)
		report.SetButton(d.DPad.Right, xbox.DPadRight)

		report.SetButton(d.Button.X, xbox.ButtonX)
		report.SetButton(d.Button.Y, xbox.ButtonY)
		report.SetButton(d.Button.A, xbox.ButtonA)
		report.SetButton(d.Button.B, xbox.ButtonB)

		report.SetButton(d.Button.Home, xbox.ButtonGuide)
		report.SetButton(d.Button.Menu, xbox.ButtonStart)
		report.SetButton(d.Button.Options, xbox.ButtonBack)

		report.SetButton(d.Stick.Left, xbox.StickLeft)
		report.SetButton(d.Stick.Right, xbox.StickRight)

		report.SetButton(d.Bumper.Left, xbox.BumperLeft)
		report.SetButton(d.Bumper.Right, xbox.BumperRight)

		report.SetStick(false, int16(d.Stick.Axis.Left.X), int16(d.Stick.Axis.Left.Y))
		report.SetStick(true, int16(d.Stick.Axis.Right.X), int16(d.Stick.Axis.Right.Y))

		report.SetTrigger(false, d.Trigger.Pressure.Left)
		report.SetTrigger(true, d.Trigger.Pressure.Right)

		con.Send(&report)
	}
}

func close() {
	stop = true
}

func msg(str string) {
	notif := toast.Notification{
		AppID:   "stadia2xbox",
		Title:   "stadia2xbox",
		Message: str,
	}
	notif.Push()
}
