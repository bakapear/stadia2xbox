package stadia

type Controller struct {
	DPad
	Button
	Bumper
	Trigger
	Stick
}

type DPad struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

type Button struct {
	A         bool
	B         bool
	X         bool
	Y         bool
	Capture   bool
	Assistant bool
	Home      bool
	Menu      bool
	Options   bool
}

type Bumper struct {
	Left  bool
	Right bool
}

type Trigger struct {
	Left  bool
	Right bool
	Pressure
}

type Pressure struct {
	Left  byte
	Right byte
}

type Stick struct {
	Left  bool
	Right bool
	Axis
}

type Axis struct {
	Left  Cord
	Right Cord
}

type Cord struct {
	X int32
	Y int32
}
