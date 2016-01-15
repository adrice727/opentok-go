package opentok

import (
	"crypto/hmac"
	"crypto/sha1"
	"github.com/adrice727/opentok/session"
	"github.com/google/go-querystring/query"
	"math/rand"
	"net/http"
	"time"
)

const (
	apiURL                = "https://api.opentok.com"
	createSessionEndpoint = "/session/create"
	apiKey                = "45456032"
	apiSecret             = "ef525d95c5c9ae83acc75bb24181e3179066d413"
	tokenSentinel         = "T!=="
)

// Opentok exposes the OpenTok API
type Opentok struct {
	apiKey    string
	apiSecret string
}

func (ot *Opentok) createSession() {

	var c = &http.Client

}

type tokenOpts struct {
	createTime uint64
	expireTime uint64
	nonce      float64 // Random number
	role       string
}

func nonce() string {
	now := int64(math.Pow10(2)) + (time.Now().UnixNano() / 10000)
	token := strconv.FormatInt(now, 10)
	return token
}

func signString(unsigned string, key string) string {
	hash := hmac.New(sha1.New, key)
	decoded :=
}

func (ot *Opentok) encodeToken(sessionID string, options tokenOpts) (token string) {
	// Seconds from epoch
	now := time.Now().Unix()

	type tokenConfig struct {
		sessionID string
		tokenOpts
	}

	if !options {
		config := &tokenConfig{sessionID, tokenOpts{now, now + (60 * 60 * 24), nonce(), "publisher"}}
	} else {
		config := &tokenConfig{sessionID, options}
	}

	v, _ := query.Values(config)
	dataString := v.Encode()
	sig := hmac.New()

}
