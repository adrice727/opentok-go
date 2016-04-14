package opentok

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	archiveEndoint = "/archive"
)

// Archive represents on OpenTok archive
type Archive struct {
	options string
}

// ArchiveOptions encapsulate the options available for archiving
// Name:        The name of the archive (for your own identification)
// HasAudio:    Whether the archive will record audio (true, the default) or
//              not (false). If you set both hasAudio and hasVideo to false, the
//              call to this method results in an error.
// HasVideo:    Whether the archive will record video (true, the default) or
//              not (false). If you set both hasAudio and hasVideo to false,
//              the call to this method results in an error.
// OutputMode:  Whether all streams in the archive are recorded to a single file
//              ("composed", the default) or to individual files ("individual").
type ArchiveOptions struct {
	Name       string
	HasAudio   bool
	HasVideo   bool
	OutputMode string
}

// StartArchive starts archiving an OpenTok session
func (ot *Opentok) StartArchive(sessionID string, options ...ArchiveOptions) Archive {

	defaultOptions := &ArchiveOptions{"", true, true, "composed"}

	// Extend returns an empty interface.  We use type assetion to convert if back to a TokenConfig
	var finalOptions interface{}

	if len(options) > 0 {
		finalOptions = Extend(defaultOptions, &ArchiveOptions{})
	} else {
		finalOptions = *defaultOptions
	}

	myOptions := finalOptions.(ArchiveOptions)

	archiveConfig := struct {
		SessionID  string `json:"sessionId"`
		Name       string `json:"name"`
		HasAudio   bool   `json:"hasAudio"`
		HasVideo   bool   `json:"hasVideo"`
		OutputMode string `json:"outputMode"`
	}{sessionID, myOptions.Name, myOptions.HasAudio, myOptions.HasVideo, myOptions.OutputMode}

	client := &http.Client{}

	endpoint := apiURL + archiveEndoint

	body, _ := json.Marshal(archiveConfig)

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "OpenTok-Go-SDK/"+version)
	req.Header.Set("X-TB-PARTNER-AUTH", ot.APIKey+":"+ot.APISecret)
	req.Header.Set("Accept", "application/json")

	a := Archive{}
	return a
}

// StopArchive starts archiving an OpenTok session
func (ot *Opentok) StopArchive() Archive {
	a := Archive{}
	return a
}
