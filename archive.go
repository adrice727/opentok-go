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
func (ot *Opentok) StartArchive(sessionID string, options ...ArchiveOptions) (*Archive, error) {

	archiveOptions := &ArchiveOptions{"OpenTok Archive", true, true, "composed"}

	if len(options) > 0 {
		Update(archiveOptions, &options[0])
	}

	archiveConfig := struct {
		SessionID  string `json:"sessionId"`
		Name       string `json:"name"`
		HasAudio   bool   `json:"hasAudio"`
		HasVideo   bool   `json:"hasVideo"`
		OutputMode string `json:"outputMode"`
	}{sessionID, archiveOptions.Name, archiveOptions.HasAudio, archiveOptions.HasVideo, archiveOptions.OutputMode}

	client := &http.Client{}

	endpoint := fmt.Sprintf("%s/v2/partner/%s/archive", apiURL, ot.APIKey)

	body := bytes.NewBufferString("")
	json.NewEncoder(body).Encode(archiveConfig)

	startArchiveError := "OT: Unable to start session archive"
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Annotate(err, startArchiveError)
	}
	ot.CommonHeaders(&req.Header)
	req.Header.Set("Content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Annotate(err, startArchiveError)
	}

	defer res.Body.Close()
	archiveData := &Archive{}
	json.NewDecoder(res.Body).Decode(archiveData)

	if err != nil {
		return nil, errors.Annotate(err, "OT: An error occurred decoding the OpenTok archive")
	}

	return archiveData, nil
}

// StopArchive starts archiving an OpenTok session
func (ot *Opentok) StopArchive(archiveID string) error {

	endpoint := fmt.Sprintf("%s/v2/partner/%s/archive/%s/stop", apiURL, ot.APIKey, archiveID)
	stopArchiveError := "OT: An error occurred while trying to stop the archive"
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.Fatal(err)
		return errors.Annotate(err, stopArchiveError)
	}

	ot.CommonHeaders(&req.Header)
	req.Header.Set("Content-type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
		return errors.Annotate(err, stopArchiveError)
	}
	return nil
}

// GetArchive return an OpenTok Archive
func (ot *Opentok) GetArchive(archiveID string) (*Archive, error) {

	endpoint := fmt.Sprintf("%s/v2/partner/%s/archive/%s", apiURL, ot.APIKey, archiveID)
	getArchiveError := "OT: An error occurred while trying to retrive the archive"
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Annotate(err, getArchiveError)
	}

	ot.CommonHeaders(&req.Header)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, errors.Annotate(err, getArchiveError)
	}

	archiveData := &Archive{}
	json.NewDecoder(res.Body).Decode(archiveData)

	return archiveData, nil
}
