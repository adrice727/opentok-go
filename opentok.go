package opentok

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"hash"
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
	req.Header.Set("User-Agent", "OpenTok-Node-SDK/"+version)
	req.Header.Set("X-TB-PARTNER-AUTH", ot.APIKey+":"+ot.APISecret)
	req.Header.Set("Accept", "application/json")

	res, _ := client.Do(req)
	decoder := json.NewDecoder(res.Body)

	var session Session
	err = decoder.Decode(&session)
	if err != nil {
		panic(err)
	}

	return session

}

func (ot *Opentok) generateToken(sessionID string) string {

	return sessionID

}

type tokenOpts struct {
	createTime uint64
	expireTime uint64
	nonce      string // Random number
	role       string
}

func nonce() string {
	now := int64(math.Pow10(2)) + (time.Now().UnixNano() / 10000)
	token := strconv.FormatInt(now, 10)
	return token
}

func signString(unsigned, key []byte) hash.Hash {
	things := hmac.New(sha1.New, append(key, unsigned...))
	return things
}

func (ot *Opentok) encodeToken(sessionID string, options ...tokenOpts) (token string) {
	// Seconds from epoch
	now := time.Now().Unix()

	type tokenConfig struct {
		sessionID string
		tokenOpts
	}

	config := &tokenConfig{sessionID, tokenOpts{uint64(now), uint64(now) + (60 * 60 * 24), nonce(), "publisher"}}

	// if len(options) == 1 {
	// 	config = &tokenConfig{sessionID, tokenOpts{uint64(now), uint64(now) + (60 * 60 * 24), nonce(), "publisher"}}
	// } else {
	// 	config = &tokenConfig{sessionID, options[0]}
	// }

	v, _ := query.Values(config)
	dataString := v.Encode()
	sig := signString([]byte(dataString), []byte(ot.APISecret))

	var decoded bytes.Buffer
	s := strings.Join([]string{"partner_id=", ot.APIKey, "&sig=", string(sig.Sum(nil)), ":", dataString}, "")
	decoded.Write([]byte(s))
	return
}
