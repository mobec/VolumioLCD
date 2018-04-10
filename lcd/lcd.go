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

	lcd.writeIR([]byte{0x03, 0x03, 0x03, 0x02})
	lcd.writeIR([]byte{
		// the PCF8574 lcd backpack has only 4 data bus lines (DB4 to DB7)
		cmdFunctionSet | func2Line | func5x8Dots | func4BitMode,
		cmdDisplayControl | displayOn,
		cmdClearDisplay,
		cmdEntryModeSet | entryLeft,
	})
	time.Sleep(time.Duration(200) * time.Millisecond)

	lcd.writeIR([]byte{backlightOn})

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

	l.writeIR([]byte{0x80 + addr})
	l.writeDR([]byte("Hello World"))

	return nil
}

//Close must be called to free underlying ressources of the LCD
func (l *LCD) Close() {
	l.dev.Close()
}

// Driver functions

//nibble splits 8bit data into nibbled 4b data + 4b signal
// data of double length
func nibble(mode byte, data []byte) []byte {
	nibBuf := make([]byte, 4*len(data))
	for i := range data {
		higher := (data[i] & 0xF0)
		lower := ((data[i] << 4) & 0xF0)
		nibBuf[i] = higher | mode
		nibBuf[i+1] = higher | mode | en
		nibBuf[i+2] = lower | mode
		nibBuf[i+3] = lower | mode | en
	}
	return nibBuf
}

//unnibble merges nibbled (4bit) data into 8bit data
func unnibble(nibBuf []byte) []byte {
	data := make([]byte, len(nibBuf)/4)
	for i := range data {
		higher := nibBuf[i] & 0xF0
		lower := (nibBuf[i+2] & 0xF0) >> 4
		data[i] = higher | lower
	}
	return data
}

func (l *LCD) writeIR(cmds []byte) error {
	// IR write as an internal operation (display clear, etc.)
	data := nibble(0x00, cmds)
	return l.dev.Write(data)
}

func (l *LCD) readIR(length int) ([]byte, error) {
	// Read busy flag (DB7) and address counter (DB0 to DB6)
	buf := nibble(rw, make([]byte, length))
	err := l.dev.Read(buf)
	return unnibble(buf), err
}

func (l *LCD) writeDR(data []byte) error {
	// DR write as an internal operation (DR to DDRAM or CGRAM
	buf := nibble(rs, data)
	return l.dev.Write(buf)
}

func (l *LCD) readDR(length int) ([]byte, error) {
	// DR read as an internal operation (DDRAM or CGRAM to DR)
	buf := nibble(rw|rs, make([]byte, length))
	err := l.dev.Read(buf)
	return unnibble(buf), err
}
