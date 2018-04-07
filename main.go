package main

// import "fmt"
import (
	"VolumioLCD/volumio"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	maxIdleConnections int = 20
	timeout            int = 10
	updateInterval     int = 200
)

func main() {
	// Create the http client
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Initialize volumio client
	volumio.URI = "http://volumio.local:3000"
	volumio.HttpClient = &httpClient

	for true {
		state, err := volumio.GetPlayerState()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\r%s (%s)", state.Title, state.Artist)

		time.Sleep(time.Duration(updateInterval) * time.Millisecond)
	}
}
