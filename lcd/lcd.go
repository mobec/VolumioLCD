package lcd

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
)

type flag byte

const (
	//commands
	cmdClearDisplay   byte = 0x01
	cmdReturnHome     byte = 0x02
	cmdEntryModeSet   byte = 0x04
	cmdDisplayControl byte = 0x08
	cmdCursorShift    byte = 0x10
	cmdFunctionSet    byte = 0x20
	cmdSetCGRAMAddr   byte = 0x40
	cmdSetDRAMAddr    byte = 0x80
	//flags for display entry mode
	entryRight          byte = 0x00
	entryLeft           byte = 0x02
	entryShiftIncrement byte = 0x01
	entryShiftDecrement byte = 0x00
	//flags for display on/off control
	displayOn        byte = 0x04
	displayOff       byte = 0x00
	displayCursorOn  byte = 0x02
	displayCursorOff byte = 0x00
	displayBlinkOn   byte = 0x01
	displayBlinkOff  byte = 0x00
	//flags for display/ cursor shift
	moveDisplay byte = 0x08
	moveCursor  byte = 0x00
	moveRight   byte = 0x04
	moveLeft    byte = 0x00
	//flags for function set
	func8BitMode byte = 0x10
	func4BitMode byte = 0x00
	func2Line    byte = 0x08
	func1Line    byte = 0x00
	func5x10Dots byte = 0x04
	func5x8Dots  byte = 0x00
	//flags for backlight
	backlightOn  byte = 0x08
	backlightOff byte = 0x00
	//other flags
	en byte = 0x04 // enable bit
	rw byte = 0x02 // read/write bit
	rs byte = 0x01 // register select bit
)

type LCD struct {
	dev *i2c.Device
}

//New opens a connection to an lcd display and sets it up
func New(line int, addr int) (*LCD, error) {
	var lcd LCD
	var err error
	lcd.dev, err = i2c.Open(&i2c.Devfs{Dev: fmt.Sprintf("/dev/i2c-%d", line)}, addr)
	if err != nil {
		return nil, err
	}

	lcd.dev.Write([]byte{0x03, 0x03, 0x03, 0x02})
	lcd.dev.Write([]byte{
		cmdFunctionSet | func2Line | func5x8Dots | func8BitMode,
		cmdDisplayControl | displayOn,
		cmdClearDisplay,
		cmdEntryModeSet | entryLeft,
	})
	time.Sleep(time.Duration(200) * time.Millisecond)
    lcd.dev.Write([]byte{backlightOn})

	return &lcd, err
}

//Show presents a string on the display
func (l *LCD) Show(str string, line uint8, pos uint8) error {
	var addr byte
	switch line {
	case 1:
		addr = byte(pos)
	case 2:
		addr = 0x40 + byte(pos)
	case 3:
		addr = 0x14 + byte(pos)
	case 4:
		addr = 0x54 + byte(pos)
	default:
		return fmt.Errorf("Line %d is not valid", line)
	}

	l.dev.Write([]byte{0x80 + addr})

	return nil
}

//Close must be called to free underlying ressources of the LCD
func (l *LCD) Close() {
	l.dev.Close()
}
