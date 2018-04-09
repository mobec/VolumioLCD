package display

import (
	"VolumioLCD/logger"
	"time"

	"github.com/davecheney/i2c"
)

type Display interface {
}

type LCD struct {
	connection *i2c.I2C
	lcd        *i2c.Lcd
	Screen     Screen
	loopStart  time.Time
	frequency  float64
}

// NewLCD given the i2c bus line and the LCD's address. You can use i2cdetect to get the address of the LCD
func NewLCD(line int, address uint8) LCD {
	var lcd LCD
	var err error
	lcd.connection, err = i2c.New(address, line)
	if err != nil {
		logger.Errorf(err.Error())
	}

	lcd.lcd, err = i2c.NewLcd(lcd.connection, 0x04, 0x01, 0x02, 0x10, 0x20, 0x40, 0x80, 0x08)
	if err != nil {
		logger.Errorf(err.Error())
	}
	lcd.lcd.BacklightOff()
	lcd.lcd.BacklightOn()

	lcd.frequency = 0.5
	lcd.Screen = NewScreen(2, 16)
	//lcd.lcd.ShowMessage("Hello World", hd44780.SHOW_LINE_1)
	//time.Sleep(time.Duration(10.0) * time.Second)

	go lcd.loop()
	return lcd
}

// Close must be called to close the connection to the lcd in a clean way
func (lcd *LCD) Close() {
	lcd.connection.Close()
}

func (lcd *LCD) loop() {
	for {
		deltaTime := time.Since(lcd.loopStart)
		lcd.loopStart = time.Now()
		// update dynamic view elements
		lcd.Screen.update(deltaTime.Seconds())
		// retrieve content from rows

		for idx := 0; idx < len(lcd.Screen.rows); idx++ {
			row := lcd.Screen.rows[idx].content()
			println(row)
			lcd.lcd.SetPosition(byte(idx+1), 0)
			size, err := lcd.lcd.Write([]byte(row))
			if err != nil {
				logger.Errorf(err.Error())
			}
			print(size)
		}
		println()
		//sleep thread to limit frequency
		time.Sleep(time.Duration(1.0/lcd.frequency)*time.Second - time.Since(lcd.loopStart))
	}
}
