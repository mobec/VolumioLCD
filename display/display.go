package display

import (
	"VolumioLCD/logger"
	"log"
	"time"

	hd44780 "github.com/d2r2/go-hd44780"
	i2c "github.com/d2r2/go-i2c"
)

type Display interface {
}

type LCD struct {
	connection *i2c.I2C
	lcd        *hd44780.Lcd
	Screen     Screen
	loopStart  time.Time
	frequency  float64
}

// NewLCD given the i2c bus line and the LCD's address. You can use i2cdetect to get the address of the LCD
func NewLCD(line int, address uint8) LCD {
	var lcd LCD
	var err error
	lcd.connection, err = i2c.NewI2C(address, line)
	if err != nil {
		log.Fatal(err)
	}

	lcd.lcd, err = hd44780.NewLcd(lcd.connection, hd44780.LCD_UNKNOWN)
	if err != nil {
		log.Fatal(err)
	}

	lcd.Screen = NewScreen(2, 16)

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

		for idx := range lcd.Screen.rows {
			row := lcd.Screen.rows[idx].content()
			if err := lcd.lcd.ShowMessage(row, hd44780.ShowOptions(idx+1)); err != nil {
				logger.Errorf(err.Error())
			}
		}
		//sleep thread to limit frequency
		time.Sleep(time.Duration(1.0/lcd.frequency)*time.Second - time.Since(lcd.loopStart))
	}
}
