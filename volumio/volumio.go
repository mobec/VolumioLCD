package volumio

import (
	"encoding/json"
	"net/http"
)

var (
	HttpClient *http.Client
	URI        string
)

const (
	maxIdleConnections int = 20
	requestTimeout     int = 10
)

func init() {
}

// PlayerState contains all the fields from a volumio REST API 'getstate' call
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

// GetPlayerState retrieves player info by calling volumios REST api. The URI of the volumio server must be provided
func GetPlayerState() (PlayerState, error) {
	var state PlayerState

	// create the request
	// FUTURE WORK: reuse request
	reqURI := URI + "/api/v1/getstate"
	req, err := http.NewRequest("GET", reqURI, nil)
	if err != nil {
		return state, err
	}

	// send the request with the reused http connection
	resp, err := HttpClient.Do(req)
	if err != nil {
		return state, err
	}
	defer resp.Body.Close()

	// decode the json response
	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return state, err
	}
	return state, nil
}
