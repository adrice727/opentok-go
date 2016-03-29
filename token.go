package opentok

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	tokenSentinel = "T1=="
)

// TokenOpts Structs tags are used by querystring package
type TokenOpts struct {
	CreateTime uint64 `url:"create_time"`
	ExpireTime uint64 `url:"expire_time"`
	Nonce      string `url:"nonce"` // Random number
	Role       string `url:"role"`
}

// TokenConfig does things
type TokenConfig struct {
	SessionID string
	Options   TokenOpts
}

// not currently used
func random() string {
	// rand should be seeded, or sourced at the top leve of the file
	// which I can do once I move all of the token functions to a
	// separate module
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	token := random.Float64()
	return strconv.FormatFloat(token, 'f', 17, 64)
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

// EncodeToken requires a tokenConfig, apiKey, and apiSecret
func EncodeToken(config TokenConfig, apiKey string, apiSecret string) string {

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
