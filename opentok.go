package opentok

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"github.com/google/go-querystring/query"
	"hash"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	apiURL                = "https://api.opentok.com"
	createSessionEndpoint = "/session/create"
	// apiKey                = "45456032"
	// apiSecret             = "ef525d95c5c9ae83acc75bb24181e3179066d413"
	tokenSentinel         = "T!=="
)

// Opentok exposes the OpenTok API
type Opentok struct {
	apiKey    string
	apiSecret string
}

func (ot *Opentok) createSession() {

	// var c = &http.Client

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
	sig := signString([]byte(dataString), []byte(apiSecret))

	var decoded bytes.Buffer
	s := strings.Join([]string{"partner_id=", apiKey, "&sig=", string(sig.Sum(nil)), ":", dataString}, "")
	decoded.Write([]byte(s))
	return
}
