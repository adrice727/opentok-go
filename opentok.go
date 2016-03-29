package opentok

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	tokenSentinel         = "T!=="
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
func (ot *Opentok) CreateSession() Session {

	client := &http.Client{}

	endpoint := apiURL + createSessionEndpoint

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "OpenTok-Go-SDK/"+version)
	req.Header.Set("X-TB-PARTNER-AUTH", ot.APIKey+":"+ot.APISecret)
	req.Header.Set("Accept", "application/json")

	res, _ := client.Do(req)
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	var sessionData []Session

	err = decoder.Decode(&sessionData)
	if err != nil {
		panic(err)
	}

	return sessionData[0]
}

// Structs tags are used by querystring package
type tokenOpts struct {
	CreateTime uint64 `url:"create_time"`
	ExpireTime uint64 `url:"expire_time"`
	Nonce      string `url:"nonce"` // Random number
	Role       string `url:"role"`
}

type tokenConfig struct {
	SessionID string
	Options   tokenOpts
}

// GenerateToken returns an opentok token
func (ot *Opentok) GenerateToken(sessionID string, options ...tokenOpts) string {

	// Seconds from epoch
	now := time.Now().Unix()

	defaultConfig := &tokenConfig{sessionID, tokenOpts{uint64(now), uint64(now) + (60 * 60 * 24), nonce(), "publisher"}}

	// Extend returns an empty interface.  We use type assetion to convert if back to a tokenConfig
	var finalConfig interface{}

	if len(options) > 0 {
		finalConfig = Extend(defaultConfig, &tokenConfig{SessionID: sessionID, Options: options[0]})
	} else {
		finalConfig = *defaultConfig
	}

	return encodeToken(finalConfig.(tokenConfig), ot.APIKey, ot.APISecret)

}

func nonce() string {
	now := int64(math.Pow10(2)) + (time.Now().UnixNano() / 10000)
	token := strconv.FormatInt(now, 10)
	return token
}

func signString(unsigned, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(unsigned))
	return hex.EncodeToString(h.Sum(nil))
}

// encodeToken requires a tokenConfig, apiKey, and apiSecret
func encodeToken(config tokenConfig, apiKey string, apiSecret string) string {

	session := struct {
		SessionID string `url:"session_id"`
	}{config.SessionID}

	s, _ := query.Values(session)
	o, _ := query.Values(config.Options)

	dataString := s.Encode() + "&" + o.Encode()

	sig := signString(dataString, apiSecret)
	queryString := strings.Join([]string{"partner_id=", apiKey, "&sig=", sig, ":", dataString}, "")
	return tokenSentinel + base64.StdEncoding.EncodeToString([]byte(queryString))
}
