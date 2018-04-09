package main

// import "fmt"
import (
	"VolumioLCD/logger"
	"VolumioLCD/volumio"
	"os"
	"os/signal"
)

const (
	updateInterval int    = 200
	volumioURI     string = "http://localhost:3000"
)

func main() {

	// Initialize volumio client
	volumio.URI = volumioURI

	// lcd := display.NewLCD(1, 0x27)
	// var artistText display.TextView
	// var titleText display.TextView
	// var titleScroll display.ScrollView
	// titleScroll.SetChild(&titleText)
	// lcd.Screen.GetRow(0).SetChild(&artistText)
	// lcd.Screen.GetRow(1).SetChild(&titleScroll)
	// defer lcd.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		for sig := range interrupt {
			print(sig)
			return
		}

		state, err := volumio.GetPlayerState()
		if err != nil {
			logger.Errorf(err.Error())
		}
		print(state.Artist)
		// artistText.SetText(state.Artist)
		// titleText.SetText(state.Title)
		//time.Sleep(time.Duration(updateInterval) * time.Millisecond)
	}

}
