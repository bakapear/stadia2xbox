package stadia

func parse(buf []byte) Controller {
	var controller Controller

	x := buf[1]
	if x != 8 {
		if x == 7 || x <= 1 {
			controller.DPad.Up = true
		} else if x >= 3 && x <= 5 {
			controller.DPad.Down = true
		}
		if x >= 5 && x <= 7 {
			controller.DPad.Left = true
		} else if x >= 1 && x <= 3 {
			controller.DPad.Right = true
		}
	}

	x = buf[2]
	if x != 0 {
		controller.Button.Capture = flag(x, 0)
		controller.Button.Assistant = flag(x, 1)
		controller.Trigger.Left = flag(x, 2)
		controller.Trigger.Right = flag(x, 3)
		controller.Button.Home = flag(x, 4)
		controller.Button.Menu = flag(x, 5)
		controller.Button.Options = flag(x, 6)
		controller.Stick.Right = flag(x, 7)
	}

	x = buf[3]
	if x != 0 {
		controller.Stick.Left = flag(x, 0)
		controller.Bumper.Right = flag(x, 1)
		controller.Bumper.Left = flag(x, 2)
		controller.Button.Y = flag(x, 3)
		controller.Button.X = flag(x, 4)
		controller.Button.B = flag(x, 5)
		controller.Button.A = flag(x, 6)
	}

	for i := 4; i <= 7; i++ {
		if buf[i] <= 127 && buf[i] >= 1 {
			buf[i]--
		}
	}

	controller.Stick.Axis.Left.X = convert(buf[4]) - 0x8000
	controller.Stick.Axis.Left.Y = -convert(buf[5]) + 0x7fff
	controller.Stick.Axis.Right.X = convert(buf[6]) - 0x8000
	controller.Stick.Axis.Right.Y = -convert(buf[7]) + 0x7fff

	controller.Trigger.Pressure.Left = buf[8]
	controller.Trigger.Pressure.Right = buf[9]

	return controller
}

func flag(x, flag byte) bool {
	flag = 1 << flag
	return x&flag == flag
}

func convert(val byte) int32 {
	value := int32(val)
	value = value<<8 | ((value << 1) & 0b1111)

	if value == 0xfffe {
		return 0xffff
	}

	return value
}
