package display

import (
	"log"

	hd44780 "github.com/d2r2/go-hd44780"
	i2c "github.com/d2r2/go-i2c"
)

type Display interface {
}

type LCD struct {
	connection *i2c.I2C
	lcd        *hd44780.Lcd
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

	return lcd
}

// Close must be called to close the connection to the lcd in a clean way
func (lcd *LCD) Close() {
	lcd.connection.Close()
}

type LCDField struct {
	lcd       *LCD
	row       int
	left      int
	right     int
	scrolling bool
}

func (lcd *LCD) NewField(row int, left int, right int, scrolling bool) LCDField {
	var field LCDField
	field.lcd = lcd
	field.row = row
	field.left = left
	field.right = right
	field.scrolling = scrolling
	return field
}
