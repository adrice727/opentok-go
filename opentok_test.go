package opentok

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type Configuration struct {
	APIKey    string
	APISecret string
}

var config Configuration

func TestConfiguration(t *testing.T) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func TestSessionCreation(t *testing.T) {
	ot := Opentok{config.APIKey, config.APISecret}
	session := ot.CreateSession()
	fmt.Printf("%+v\n", session)
}

func TestTokenCreation(t *testing.T) {
	ot := Opentok{config.APIKey, config.APISecret}
	session := ot.CreateSession()
	token := ot.GenerateToken(session.SessionID)
	fmt.Println("session", session)
	fmt.Printf("%+v\n", token)
}
