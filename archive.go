package opentok

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/juju/errors"
)

const (
	archiveEndoint = "/v2/partner/archive"
)

// Archive represents on OpenTok archive
// type Archive struct {
// 	options string
// }

// Archive represents on OpenTok archive
type Archive struct {

	// Unix timestamp that specified when the
	// archive was created
	CreatedAt int64 `json:"createdAt"`

	// Duration of the archive in seconds
	Duration int64 `json:"duration"`

	// ID of the archive. This is a unique id
	// identifier for the archive. It's used to
	// stop, retrieve and delete archives
	ID string `json:"id"`

	// Name for the archive. The user can choose
	// any name but it doesn't necessarily need
	// to be different between archives
	Name string `json:"name"`

	// APIKey to which the archive belongs
	APIKey int `json:"partnerId"`

	// SessionID to which the archive belongs
	SessionID string `json:"sessionId"`

	// Size of the archives in KB
	Size int `json:"size"`

	// URL from where the archive can be retrieved. This is
	// only useful if the archive is in status available
	// in the OpenTok S3 Account
	URL string `json:"url"`

	// Status of the Archive. The possibilities are:
	// - `started`: if the archive is being recorded
	// - `stopped`: if the archive has been stopped and it hasn't
	//   been uploaded or available
	// - `deleted`: if the archive has been deleted. Only available
	//   archives can be deleted
	// - `uploaded`: if the archive has been uploaded to the
	//   partner storage account
	// - `paused`: if the archive has not been stopped but it is not
	//   recording. It can transition to Started again
	// - `available`: if the archive has been uploaded to the
	//   OpenTok S3 account
	// - `expired`: available archives are removed from the OpenTok
	//   S3 account after 3 days. Their status become expired.
	Status string `json:"status"`

	// HasAudio tells whether the archive contains an audio
	// stream.
	HasAudio bool `json:"hasAudio"`

	// HasVideo tells whether the archive contains a video
	// stream.
	HasVideo bool `json:"hasVideo"`
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
