package main

// import "fmt"
import (
	"VolumioLCD/volumio"
	"fmt"
	"log"
	"time"
)

const (
	updateInterval int    = 200
	volumioURI     string = "http://volumio.local:3000"
)

func main() {

	// Initialize volumio client
	volumio.URI = volumioURI

	for true {
		state, err := volumio.GetPlayerState()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\r%s (%s)", state.Title, state.Artist)

		time.Sleep(time.Duration(updateInterval) * time.Millisecond)
	}
}
