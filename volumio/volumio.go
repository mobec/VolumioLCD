package volumio

import (
	"encoding/json"
	"net/http"
)

type PlayerState struct {
	Status            string
	Title             string
	Artist            string
	Album             string
	Position          int
	Albumart          string
	URI               string
	Tracktype         string
	Seek              int
	Duration          int
	Samplerate        string
	Bitdepth          string
	Channels          int
	Volume            int
	Random            bool
	Repeat            bool
	Repeatsingle      bool
	Consume           bool
	Mute              bool
	Stream            bool
	Updatedb          bool
	Volatile          bool
	Service           string
	Disableuicontrols bool
}

func GetPlayerState() (PlayerState, error) {
	resp, err := http.Get("http://volumio.local:3000/api/v1/getstate")
	var state PlayerState
	if err != nil {
		return state, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return state, err
	}
	return state, nil
}
