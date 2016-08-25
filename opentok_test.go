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
	session, _ := ot.CreateSession()
	fmt.Printf("%+v\n", session)
}

func TestTokenCreation(t *testing.T) {
	ot := Opentok{config.APIKey, config.APISecret}
	session, _ := ot.CreateSession()
	token := ot.GenerateToken(session.SessionID)
	fmt.Println("session", session)
	fmt.Printf("%+v\n", token)
}

// There must be an active publisher in order to start archiving
// TODO Update error checking in Archiving
func TestStartArchive(t *testing.T) {
	ot := Opentok{config.APIKey, config.APISecret}
	session, _ := ot.CreateSession()
	archive, _ := ot.StartArchive(session.SessionID)
	fmt.Printf("Archive created for session %s", session.SessionID)
	fmt.Print(archive)
}

func TestStopArchive(t *testing.T) {
	ot := Opentok{config.APIKey, config.APISecret}
	archiveID := ""
	err := ot.StopArchive(archiveID)
	if err != nil {
		fmt.Println("Error stopping the archive")
	} else {
		fmt.Println("Succesfully stopped the archive")
	}
}

func TestGetArchive(t *testing.T) {
	ot := Opentok{config.APIKey, config.APISecret}
	archiveID := "31cd48a1-654b-4801-84dd-77c20b008a86"
	archive, err := ot.GetArchive(archiveID)
	if err != nil {
		fmt.Printf("Error retrieving archive %s", archiveID)
	} else {
		fmt.Println("Succesfully retrived the archive", archive)
	}
}

func TestDeleteArchive(t *testing.T) {
	ot := Opentok{config.APIKey, config.APISecret}
	archiveID := "31cd48a1-654b-4801-84dd-77c20b008a86"
	err := ot.DeleteArchive(archiveID)
	if err != nil {
		fmt.Printf("Error deleting archive %s", archiveID)
	} else {
		fmt.Printf("Succesfully deleted the archive %s", archiveID)
	}
}
