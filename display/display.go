package display

import (
	"VolumioLCD/lcd"
	"VolumioLCD/logger"
	"time"
)

type IDisplay interface {
}

type Display struct {
	lcd       *lcd.LCD
	Screen    *Screen
	loopStart time.Time
	Backlight bool // access is atomic, no synchronisation required
	frequency float64
}

// NewLCD given the i2c bus line and the LCD's address. You can use i2cdetect to get the address of the LCD
func New(line int, address int) *Display {
	var d Display
	var err error
	d.lcd, err = lcd.New(line, address)
	if err != nil {
		logger.Errorf(err.Error())
	}

	d.frequency = 10
	d.Screen = NewScreen(2, 16)
	//time.Sleep(time.Duration(10.0) * time.Second)

	go d.loop()
	return &d
}

// Close must be called to close the connection to the lcd in a clean way
func (d *Display) Close() {
	d.lcd.Close()
}

func (d *Display) loop() {
	for {
		deltaTime := time.Since(d.loopStart)
		d.loopStart = time.Now()
		// update dynamic view elements
		d.Screen.update(deltaTime.Seconds())
		// retrieve content from rows

		for idx := 0; idx < len(d.Screen.rows); idx++ {
			row := d.Screen.rows[idx].content()
			err := d.lcd.Show(row, uint8(idx+1), 0)
			if err != nil {
				logger.Errorf(err.Error())
			}
		}

		d.lcd.Backlight(d.Backlight)

		//sleep thread to limit frequency
		time.Sleep(time.Duration(1.0/d.frequency)*time.Second - time.Since(d.loopStart)*time.Second)
	}
}
