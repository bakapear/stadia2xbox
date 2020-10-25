# stadia2xbox
Emulates a Xbox 360 Controller from a Stadia Controller

Originally a fork of https://github.com/71/stadiacontroller/ but heavily rewritten.

### Installation
1. Install [ViGEm](https://github.com/ViGEm/ViGEmBus/releases)
2. Download the latest release from the [releases page](https://github.com/bakapear/stadia2xbox/releases)
3. Run stadia2xbox.exe

### Troubleshooting
The program tries to open your stadia controller in exclusive mode which requires it to not be accessed by any other program. To regain access try the following:
- Close any application that might use the controller (Steam, Browser, etc...)
- Replug the controller
- Go to the device manager and disable/enable "HID-compliant game controller"
- Restart computer
