package main

import (
    "VolumioLCD/display"
    "VolumioLCD/logger"
    "VolumioLCD/volumio"
    "VolumioLCD/utils"
    "context"
    "os"
    "os/signal"
    "time"
)

// import "fmt"

const (
    updateInterval int    = 200
    volumioURI     string = "http://localhost:3000"
)

func main() {
    // Correctly handle os messages
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

    timeout := 60.0 * time.Second;
    for(!utils.PathExists("/dev/i2c-1") && timeout > 0.0) {
        waitTime := 1.0 * time.Second
        time.Sleep(waitTime)
        timeout -= waitTime
    }
    time.Sleep(5.0 * time.Second)
    lcd := display.New(1, 0x27)
    lcd.Screen = display.NewScreen(2, 16)
    var artistText display.TextView
    var titleText display.TextView
    var titleScroll display.ScrollView
    titleScroll.SetChild(&titleText)
    titleScroll.SetLength(16)
    titleScroll.SetSpeed(2)
    lcd.Screen.GetRow(0).SetChild(&artistText)
    lcd.Screen.GetRow(1).SetChild(&titleScroll)
    
    go lcd.Loop()
    
    defer lcd.Close()

    go func() {
        select {
        case <-sig:
            cancel()
        case <-ctx.Done():
        }
    }()

    go func() {
        // allow to check for changes
        var prevTitle string
        var prevArtist string

        // main loop
        for {
            time.Sleep(time.Duration(updateInterval) * time.Millisecond)
            state, err := volumio.GetPlayerState()
            if err != nil {
                logger.Warningf(err.Error())
                continue
            }
            artistText.SetText(state.Artist)
            titleText.SetText(state.Title)
            if state.Status == "play" &&
            (state.Artist != prevArtist || state.Title != prevTitle) &&
            state.Artist != "" &&
            state.Title != "" {
                lcd.Backlight = true
            } else if state.Status == "stop" || state.Status == "pause" {
                lcd.Backlight = false
            }
        }
    }()

    <-sig
}
