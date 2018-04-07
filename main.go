package main

// import "fmt"
import (
	"VolumioLCD/volumio"
	"log"
)

func main() {
	state, err := volumio.GetPlayerState()
	if err != nil {
		log.Fatal(err)
	}
	print(state.Service)
}
