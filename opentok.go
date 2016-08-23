package opentok

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/juju/errors"
)

const (
	apiURL                = "https://api.opentok.com"
	createSessionEndpoint = "/session/create"
	version               = "0.0.1"
)

// Opentok exposes the OpenTok API
type Opentok struct {
	APIKey    string
	APISecret string
}

// CreateSession creates a new OpenTok session
func (ot *Opentok) CreateSession() (*Session, error) {

	client := &http.Client{}

	endpoint := apiURL + createSessionEndpoint

	req, err := http.NewRequest("POST", endpoint, nil)
	// req.Close = true
	if err != nil {
		log.Fatal(err)
		return nil, errors.Annotate(err, "OT: Unable to create an OpenTok session")

	}

	ot.CommonHeaders(&req.Header)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Annotate(err, "OT: Unable to create a session")
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	var sessionData []Session

	err = decoder.Decode(&sessionData)
	if err != nil {
		return nil, errors.Annotate(err, "OT: An error occurred decoding the OpenTok session")
	}

	return &sessionData[0], nil
}

// GenerateToken returns an opentok token
func (ot *Opentok) GenerateToken(sessionID string, options ...TokenOptions) string {

	// Seconds from epoch
	now := time.Now().Unix()

	// Default configuration
	config := &TokenConfig{sessionID, TokenOptions{uint64(now), uint64(now) + (60 * 60 * 24), nonce(), "publisher"}}

	// Extend default config with passed in options
	if len(options) > 0 {
		Update(config, &TokenConfig{SessionID: sessionID, Options: options[0]})
	}

	return EncodeToken(*config, ot.APIKey, ot.APISecret)

}

// CommonHeaders sets the common headers for requests to the OpenTok API
func (ot *Opentok) CommonHeaders(h *http.Header) {
	h.Add("User-Agent", "OpenTok-Go-SDK/"+version)
	h.Add("X-TB-PARTNER-AUTH", ot.APIKey+":"+ot.APISecret)
	h.Add("X-TB-VERSION", "1")
	h.Add("Accept", "application/json")
	h.Add("Content-type", "application/json")
}
