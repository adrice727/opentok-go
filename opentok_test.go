package opentok

import (
	"fmt"
	"testing"
)

const (
	apiKey    = "45456032"
	apiSecret = "ef525d95c5c9ae83acc75bb24181e3179066d413"
)

func TestSessionCreation(t *testing.T) {
	ot := Opentok{apiKey, apiSecret}
	session := ot.CreateSession()
	fmt.Printf("%+v\n", session)
}

func TestTokenCreation(t *testing.T) {
	ot := Opentok{apiKey, apiSecret}
	session := ot.CreateSession()
	token := ot.GenerateToken(session.SessionID)
	fmt.Printf("%+v\n", token)
}
