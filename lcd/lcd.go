package lcd

import (
	"fmt"
	"sync"
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
	// mutex protecting the ressource
	mutex sync.Mutex
	dev   *i2c.Device
	bl    byte
}

//New opens a connection to an lcd display and sets it up
func New(line int, addr int) (*LCD, error) {
	var lcd LCD
	var err error
	lcd.dev, err = i2c.Open(&i2c.Devfs{Dev: fmt.Sprintf("/dev/i2c-%d", line)}, addr)
	if err != nil {
		return nil, err
	}

	// Initialization sequence
	lcd.initialize()
	time.Sleep(time.Duration(200) * time.Millisecond)
	lcd.writeIR(cmdDisplayControl | displayOn)

	//lcd.dev.Write([]byte{backlightOn})
	lcd.bl = backlightOn
	lcd.writeIR(0x00)

	return &lcd, err
}

// initialize the HD44780U lcd with initialization by instruction
func (l *LCD) initialize() error {
	// write 0 0 1 1
	err := l.writeToDev([]byte{0x03})
	if err != nil {
		return err
	}
	// wait for more than 4.1 ms
	time.Sleep(time.Duration(5) * time.Millisecond)
	// write 0 0 1 1
	err = l.writeToDev([]byte{0x03})
	if err != nil {
		return err
	}
	// wait for more than 100 μs
	time.Sleep(time.Duration(150) * time.Microsecond)
	err = l.writeToDev([]byte{0x03})
	if err != nil {
		return err
	}

	// function set to 4 bit interface
	err = l.writeToDev([]byte{0x02})
	if err != nil {
		return err
	}

	// from here on we are in 4 bit mode and we need to nibble

	// function set the number of display lines and character font
	err = l.writeIR(cmdFunctionSet | func2Line | func5x8Dots | func4BitMode)
	if err != nil {
		return err
	}

	err = l.writeIR(cmdDisplayControl | displayOff)
	if err != nil {
		return err
	}

	err = l.writeIR(cmdClearDisplay)
	if err != nil {
		return err
	}

	err = l.writeIR(cmdEntryModeSet | entryLeft | entryShiftIncrement)
	if err != nil {
		return err
	}

	return nil
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

	l.mutex.Lock()
	defer l.mutex.Unlock()
	if err := l.writeIR(cmdSetDRAMAddr | addr); err != nil {
		return err
	}
	for i := range str {
		if err := l.writeDR(str[i]); err != nil {
			return err
		}
	}
	return nil
}

//Backlight allows to turn the lcd's backlight on and off
func (l *LCD) Backlight(isOn bool) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if isOn {
		l.bl = backlightOn
	} else {
		l.bl = backlightOff
	}
	return l.writeIR(0x00)
}

//Clear clears the display from characters
func (l *LCD) Clear() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.writeIR(cmdClearDisplay)
}

//Close must be called to free underlying ressources of the LCD
func (l *LCD) Close() {
	l.Clear()
	l.Backlight(false)

	l.mutex.Lock()
	l.dev.Close()
	l.mutex.Unlock()
}

// Driver functions

//nibble splits 8bit data into nibbled 4b data + 4b signal
// data of double length
func nibble(mode byte, data byte) []byte {
	nibBuf := make([]byte, 2)
	higher := (data & 0xF0)
	lower := ((data << 4) & 0xF0)
	nibBuf[0] = higher | mode
	nibBuf[1] = lower | mode
	return nibBuf
}

//unnibble merges nibbled (4bit) data into 8bit data
func unnibble(nibBuf []byte) byte {
	higher := nibBuf[0] & 0xF0
	lower := (nibBuf[1] & 0xF0) >> 4
	return higher | lower
}

func (l *LCD) writeIR(cmd byte) error {
	// IR write as an internal operation (display clear, etc.)
	data := nibble(0x00, cmd)
	// return l.dev.Write(data)
	return l.writeToDev(data)
}

func (l *LCD) readIR() (byte, error) {
	// Read busy flag (DB7) and address counter (DB0 to DB6)
	buf := nibble(rw, 0x00)
	err := l.dev.Read(buf)
	return unnibble(buf), err
}

func (l *LCD) writeDR(data byte) error {
	// DR write as an internal operation (DR to DDRAM or CGRAM
	buf := nibble(rs, data)
	return l.writeToDev(buf)
}

func (l *LCD) readDR() (byte, error) {
	// DR read as an internal operation (DDRAM or CGRAM to DR)
	buf := nibble(rw|rs, 0x00)
	err := l.dev.Read(buf)
	return unnibble(buf), err
}

func (l *LCD) writeToDev(data []byte) error {
	strobeBuf := make([]byte, 3*len(data))
	for i := range data {
		strobeBuf[i*3] = data[i] | l.bl
		strobeBuf[i*3+1] = data[i] | en | l.bl
		strobeBuf[i*3+2] = data[i] | l.bl
	}
	return l.dev.Write(strobeBuf)
}
