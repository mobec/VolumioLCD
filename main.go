package main

// import "fmt"
import (
	"VolumioLCD/display"
	"VolumioLCD/logger"
	"VolumioLCD/volumio"
	"time"
)

const (
	updateInterval int    = 200
	volumioURI     string = "http://localhost:3000"
)

func main() {

	// Initialize volumio client
	volumio.URI = volumioURI

	lcd := display.NewLCD(1, 0x27)
	var artistText display.TextView
	var titleText display.TextView
	var titleScroll display.ScrollView
	titleScroll.SetChild(&titleText)
	lcd.Screen.GetRow(0).SetChild(&artistText)
	lcd.Screen.GetRow(1).SetChild(&titleScroll)
	defer lcd.Close()

	for {
		state, err := volumio.GetPlayerState()
		if err != nil {
			logger.Errorf(err.Error())
		}
		println(state.Artist)
		artistText.SetText(state.Artist)
		titleText.SetText(state.Title)
		time.Sleep(time.Duration(updateInterval) * time.Millisecond)
	}

}
