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
	if err != nil {
		log.Fatal(err)
		return nil, errors.Annotate(err, "OT: Unable to create an OpenTok session")

	}
	req.Header.Set("User-Agent", "OpenTok-Go-SDK/"+version)
	req.Header.Set("X-TB-PARTNER-AUTH", ot.APIKey+":"+ot.APISecret)
	req.Header.Set("Accept", "application/json")

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
