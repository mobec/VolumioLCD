package main

import (
	"VolumioLCD/display"
	"VolumioLCD/logger"
	"VolumioLCD/volumio"
	"context"
	"os"
	"os/signal"
	"time"
)

// import "fmt"

const (
	updateInterval int    = 500
	volumioURI     string = "http://volumio.local:3000"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer func() {
		signal.Stop(sig)
		cancel()
	}()

	// Initialize volumio client
	volumio.URI = volumioURI

	lcd := display.New(1, 0x27)
	lcd.Screen = display.NewScreen(2, 16)
	var artistText display.TextView
	var titleText display.TextView
	var titleScroll display.ScrollView
	titleScroll.SetChild(&titleText)
	titleScroll.SetLength(16)
	lcd.Screen.GetRow(0).SetChild(&artistText)
	lcd.Screen.GetRow(1).SetChild(&titleScroll)
	defer lcd.Close()

	go func() {
		select {
		case <-sig:
			cancel()
		case <-ctx.Done():
		}
	}()

	go func() {
		for {
			state, err := volumio.GetPlayerState()
			if err != nil {
				logger.Warningf(err.Error())
			}
			println(state.Artist)
			artistText.SetText(state.Artist)
			titleText.SetText(state.Title)
			time.Sleep(time.Duration(updateInterval) * time.Millisecond)
		}
	}()

	<-sig
}

// func main() {
// 	lcd, err := lcd.New(1, 0x27)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer lcd.Close()

// 	lcd.Show("Hello World", 1, 0)

// 	time.Sleep(time.Duration(5) * time.Second)
// }
