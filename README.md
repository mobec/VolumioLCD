VolumioLCD
==========

VolumioLCD is a tool to display the Volumio player state on hd44780 LCD displays over I2C on a Raspberry Pi.

##Features
* Works with all services supported by volumio (by using volumios REST API)
* No additional tools required (due to a custom hd44780 driver for the LCD)

##Current Limitations
* Only 16x2 displays are supported
* No config, therefore hardcoded I2C address and line and volumio player URI

##Installation
1) Install golang using the official .tar (do NOT use apt-get)
2) Clone the VolumioLCD repository to $GOPATH/src
3) inside the VolumioLCD directory: ```go get`` and then ```go build```
4) start VolumioLCD on every boot (e.g. using systemd)
