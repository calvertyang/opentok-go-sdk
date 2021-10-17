package opentok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// LayoutType is the alias of string type.
type LayoutType string

const (
	// BestFit is a tiled layout, which scales according to the number of videos.
	// The number of columns and rows of varies depending on the number of
	// OpenTok streams in the archive.
	BestFit LayoutType = "bestFit"

	// PIP is a picture-in-picture layout, where a small stream is visible over
	// a full-size stream.
	PIP LayoutType = "pip"

	// VerticalPresentation is a layout with one large stream on the right edge of
	// the output, and several smaller streams along the left edge of the output.
	VerticalPresentation LayoutType = "verticalPresentation"

	// HorizontalPresentation is a layout with one large stream on the top edge of
	// the output, and several smaller streams along the bottom edge of the output.
	HorizontalPresentation LayoutType = "horizontalPresentation"

	// Custom let you can use CSS to define your own custom layout for composed
	// archive videos.
	Custom LayoutType = "custom"
)

// ArchiveOutputMode is the alias of string type.
type ArchiveOutputMode string

const (
	// Composed means the archive is a single MP4 file composed of all streams.
	Composed ArchiveOutputMode = "composed"

	// Individual means the archive is a ZIP container file with multiple
	// individual media files for each stream, and a JSON metadata file for video
	// synchronization.
	Individual ArchiveOutputMode = "individual"
)

// Resolution is the alias of string type.
type Resolution string

// The resolution of the archive, either SD(default) or HD.
const (
	// SD (640x480-pixel) archives have a 4:3 aspect ratio.
	SD Resolution = "640x480"
	// HD (1280x720-pixel) archives have a 16:9 aspect ratio.
	HD Resolution = "1280x720"
)

// Layout defines the layout type for the archive.
type Layout struct {
	Type       LayoutType `json:"type,omitempty"`
	StyleSheet string     `json:"stylesheet,omitempty"`
}

// ArchiveOptions defines the options for starting an archive recording.
type ArchiveOptions struct {
	// The session ID of the OpenTok session you want to start archiving.
	SessionID string `json:"sessionId"`

	// Whether the archive will record audio (true, the default) or not (false).
	HasAudio bool `json:"hasAudio,omitempty"`

	// Whether the archive will record video (true, the default) or not (false).
	HasVideo bool `json:"hasVideo,omitempty"`

	// Specify this to assign the initial layout type for the archive.
	Layout *Layout `json:"layout,omitempty"`

	// The name of the archive (for your own identification).
	Name string `json:"name,omitempty"`

	// Whether all streams in the archive are recorded to a single file or to
	// individual files.
	OutputMode ArchiveOutputMode `json:"outputMode,omitempty"`

	// The resolution of the archive.
	Resolution Resolution `json:"resolution,omitempty"`
}

// Archive defines the response returned from API.
type Archive struct {
	// The time at which the archive was created, in milliseconds since the
	// UNIX epoch.
	CreatedAt int `json:"createdAt"`

	// The duration of the archive, in milliseconds.
	Duration int `json:"duration"`

	// Whether the archive has an audio track or not.
	HasAudio bool `json:"hasAudio"`

	// Whether the archive has an video track or not.
	HasVideo bool `json:"hasVideo"`

	// The unique archive ID.
	ID string `json:"id"`

	// The name of the archive.
	Name string `json:"name"`

	// The output mode to be generated for this archive.
	OutputMode ArchiveOutputMode `json:"outputMode"`

	// The API key associated with the archive.
	ProjectID int `json:"projectId"`

	// This string describes the reason the archive stopped or failed.
	Reason string `json:"reason"`

	// The resolution of the archive.
	Resolution Resolution `json:"resolution"`

	// The session ID of the OpenTok session associated with this archive.
	SessionID string `json:"sessionId"`

	// The size of the MP4 file.
	Size int `json:"size"`

	// The status of the archive.
	Status string `json:"status"`

	// The download URL of the available MP4 file.
	URL *string `json:"url"`

	// The instance of OpenTok.
	OpenTok *OpenTok `json:"-"`
}

// ArchiveListOptions defines the query parameters to filter the list of
// archives.
type ArchiveListOptions struct {
	// Query parameters to specify the index offset of the first archive.
	Offset int

	// Query parameter to limit the number of archives to be returned.
	Count int

	// Query parameter to list archives for a specific session ID.
	SessionID string
}

// ArchiveList defines the response returned from API.
type ArchiveList struct {
	// The total number of archives for the API key.
	Count int `json:"count"`

	// An array of objects defining each archive retrieved.
	// Archives are listed from the newest to the oldest in the return set.
	Items []*Archive `json:"items"`
}

// AmazonS3Config defines the Amazon S3 config for setting archiving upload
// target.
type AmazonS3Config struct {
	// The Amazon Web Services access key.
	AccessKey string `json:"accessKey"`

	// The Amazon Web Services secret key.
	SecretKey string `json:"secretKey"`

	// The S3 bucket name.
	Bucket string `json:"bucket"`

	// The S3 or S3-compatible storage endpoint.
	Endpoint string `json:"endpoint,omitempty"`
}

// AzureConfig defines the Microsoft Azure config for setting archiving upload
// target.
type AzureConfig struct {
	// The Microsoft Azure account name.
	AccountName string `json:"accountName"`

	// The Microsoft Azure account key.
	AccountKey string `json:"accountKey"`

	// The Microsoft Azure container name.
	Container string `json:"container"`

	// The Microsoft Azure domain in which the container resides.
	Domain string `json:"domain,omitempty"`
}

// StorageOptions defines the options for setting archiving upload target.
type StorageOptions struct {
	// Type of upload target.
	Type string `json:"type"`

	// Settings for the target.
	Config interface{} `json:"config"`

	// Set this to "opentok" to have the archive available at the OpenTok
	// dashboard if upload fails.
	// Set this to "none" (or omit the property) to prevents archive files from
	// being stored in the OpenTok cloud if the upload fails.
	Fallback string `json:"fallback,omitempty"`
}

// StartArchive starts the recording of the archive.
//
// To successfully start recording an archive, at least one client must be
// connected to the session.
// You can only record one archive at a time for a given session.
// You can only record archives of sessions that use the OpenTok Media Router.
func (ot *OpenTok) StartArchive(sessionID string, opts *ArchiveOptions) (*Archive, error) {
	return ot.StartArchiveContext(context.Background(), sessionID, opts)
}

// StartArchiveContext uses ctx for HTTP requests.
func (ot *OpenTok) StartArchiveContext(ctx context.Context, sessionID string, opts *ArchiveOptions) (*Archive, error) {
	opts.SessionID = sessionID

	if opts.Layout != nil {
		if opts.Layout.Type != BestFit && opts.Layout.Type != PIP && opts.Layout.Type != Custom &&
			opts.Layout.Type != VerticalPresentation && opts.Layout.Type != HorizontalPresentation {
			return nil, fmt.Errorf("Invalid type of layout for start archive")
		}

		if opts.Layout.Type == Custom && opts.Layout.StyleSheet == "" {
			return nil, fmt.Errorf("StyleSheet property of layout cannot be empty")
		}

		// For other layout types, do not set a stylesheet property.
		if opts.Layout.Type != Custom && opts.Layout.StyleSheet != "" {
			return nil, fmt.Errorf("Set stylesheet property only when using custom layout")
		}
	}

	if opts.OutputMode != "" && opts.OutputMode != Composed && opts.OutputMode != Individual {
		return nil, fmt.Errorf("Invalid output mode for start archive")
	}

	if opts.Resolution != "" && opts.Resolution != SD && opts.Resolution != HD {
		return nil, fmt.Errorf("Invalid resolution for start archive")
	}

	jsonStr, _ := json.Marshal(opts)

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	archive := &Archive{}
	if err := json.NewDecoder(res.Body).Decode(archive); err != nil {
		return nil, err
	}

	archive.OpenTok = ot

	return archive, nil
}

// StopArchive stops the recording of the archive.
//
// Archives stop recording after 2 hours (120 minutes), or 60 seconds after the
// last client disconnects from the session, or 60 minutes after the last
// client stops publishing.
func (ot *OpenTok) StopArchive(archiveID string) (*Archive, error) {
	return ot.StopArchiveContext(context.Background(), archiveID)
}

// StopArchiveContext uses ctx for HTTP requests.
func (ot *OpenTok) StopArchiveContext(ctx context.Context, archiveID string) (*Archive, error) {
	if archiveID == "" {
		return nil, fmt.Errorf("Archive recording cannot be stopped without an archive ID")
	}

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive/" + archiveID + "/stop"
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	archive := &Archive{}
	if err = json.NewDecoder(res.Body).Decode(archive); err != nil {
		return nil, err
	}

	archive.OpenTok = ot

	return archive, nil
}

// ListArchives returns the records of all archives for your project that are
// in progress.
func (ot *OpenTok) ListArchives(opts *ArchiveListOptions) (*ArchiveList, error) {
	return ot.ListArchivesContext(context.Background(), opts)
}

// ListArchivesContext uses ctx for HTTP requests.
func (ot *OpenTok) ListArchivesContext(ctx context.Context, opts *ArchiveListOptions) (*ArchiveList, error) {
	params := []string{"?"}

	if opts.Offset != 0 {
		params = append(params, "offset="+strconv.Itoa(opts.Offset))
	}

	if opts.Count != 0 {
		params = append(params, "count="+strconv.Itoa(opts.Count))
	}

	if opts.SessionID != "" {
		params = append(params, "sessionId="+opts.SessionID)
	}

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive" + strings.Join(params, "&")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	archiveList := &ArchiveList{}
	if err := json.NewDecoder(res.Body).Decode(archiveList); err != nil {
		return nil, err
	}

	for _, archive := range archiveList.Items {
		archive.OpenTok = ot
	}

	return archiveList, nil
}

// GetArchive returns a archive details record describing the archive.
func (ot *OpenTok) GetArchive(archiveID string) (*Archive, error) {
	return ot.GetArchiveContext(context.Background(), archiveID)
}

// GetArchiveContext uses ctx for HTTP requests.
func (ot *OpenTok) GetArchiveContext(ctx context.Context, archiveID string) (*Archive, error) {
	if archiveID == "" {
		return nil, fmt.Errorf("Cannot get archive information without an archive ID")
	}

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive/" + archiveID
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	archive := &Archive{}
	if err := json.NewDecoder(res.Body).Decode(archive); err != nil {
		return nil, err
	}

	archive.OpenTok = ot

	return archive, nil
}

// DeleteArchive deletes a specific archive.
func (ot *OpenTok) DeleteArchive(archiveID string) error {
	return ot.DeleteArchiveContext(context.Background(), archiveID)
}

// DeleteArchiveContext uses ctx for HTTP requests.
func (ot *OpenTok) DeleteArchiveContext(ctx context.Context, archiveID string) error {
	if archiveID == "" {
		return fmt.Errorf("Archive cannot be deleted without an archive ID")
	}

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive/" + archiveID
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	return nil
}

// SetArchiveStorage let you can have OpenTok upload completed archives to an
// Amazon S3 bucket (or an S3-compliant storage provider) or Microsoft Azure container.
func (ot *OpenTok) SetArchiveStorage(opts *StorageOptions) (*StorageOptions, error) {
	return ot.SetArchiveStorageContext(context.Background(), opts)
}

// SetArchiveStorageContext uses ctx for HTTP requests.
func (ot *OpenTok) SetArchiveStorageContext(ctx context.Context, opts *StorageOptions) (*StorageOptions, error) {
	if opts.Type != "s3" && opts.Type != "azure" {
		return nil, fmt.Errorf("Only support Amazon S3 or Microsoft Azure for upload completed archives")
	}

	switch config := opts.Config.(type) {
	case *AmazonS3Config:
		if config.AccessKey == "" {
			return nil, fmt.Errorf("The Amazon Web Services access key cannot be empty")
		}

		if config.SecretKey == "" {
			return nil, fmt.Errorf("The Amazon Web Services secret key cannot be empty")
		}

		if config.Bucket == "" {
			return nil, fmt.Errorf("The S3 bucket name cannot be empty")
		}
	case *AzureConfig:
		if config.AccountName == "" {
			return nil, fmt.Errorf("The Microsoft Azure account name cannot be empty")
		}

		if config.AccountKey == "" {
			return nil, fmt.Errorf("The Microsoft Azure account key cannot be empty")
		}

		if config.Container == "" {
			return nil, fmt.Errorf("The Microsoft Azure container name cannot be empty")
		}
	default:
		return nil, fmt.Errorf("Invalid archive storage config")
	}

	jsonStr, _ := json.Marshal(opts)

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive/storage"
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	options := &StorageOptions{}
	if err := json.NewDecoder(res.Body).Decode(options); err != nil {
		return nil, err
	}

	return options, nil
}

// DeleteArchiveStorage deletes the configuration of archive storage.
func (ot *OpenTok) DeleteArchiveStorage() error {
	return ot.DeleteArchiveStorageContext(context.Background())
}

// DeleteArchiveStorageContext uses ctx for HTTP requests.
func (ot *OpenTok) DeleteArchiveStorageContext(ctx context.Context) error {
	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive/storage"
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	return nil
}

// SetArchiveLayout dynamically change the layout type of a composed archive.
func (ot *OpenTok) SetArchiveLayout(archiveID string, layout *Layout) (*Archive, error) {
	return ot.SetArchiveLayoutContext(context.Background(), archiveID, layout)
}

// SetArchiveLayoutContext uses ctx for HTTP requests.
func (ot *OpenTok) SetArchiveLayoutContext(ctx context.Context, archiveID string, layout *Layout) (*Archive, error) {
	if archiveID == "" {
		return nil, fmt.Errorf("Cannot change the layout type of a composed archive without an archive ID")
	}

	if layout.Type != BestFit && layout.Type != PIP && layout.Type != Custom &&
		layout.Type != VerticalPresentation && layout.Type != HorizontalPresentation {
		return nil, fmt.Errorf("Invalid type of layout for archive")
	}

	if layout.Type == Custom && layout.StyleSheet == "" {
		return nil, fmt.Errorf("StyleSheet property of layout cannot be empty")
	}

	// For other layout types, do not set a stylesheet property.
	if layout.Type != Custom && layout.StyleSheet != "" {
		return nil, fmt.Errorf("Set stylesheet property only when using custom layout")
	}

	jsonStr, _ := json.Marshal(layout)

	// Create jwt token
	jwt, err := ot.jwtToken(projectToken)
	if err != nil {
		return nil, err
	}

	endpoint := ot.apiHost + projectURL + "/" + ot.apiKey + "/archive/" + archiveID + "/layout"
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)
	req.Header.Add("User-Agent", SDKName+"/"+SDKVersion)

	res, err := ot.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	archive := &Archive{}
	if err := json.NewDecoder(res.Body).Decode(archive); err != nil {
		return nil, err
	}

	archive.OpenTok = ot

	return archive, nil
}

// Stop stops the recording of the archive.
func (archive *Archive) Stop() (*Archive, error) {
	return archive.OpenTok.StopArchive(archive.ID)
}

// Delete deletes a specific archive.
func (archive *Archive) Delete() error {
	return archive.OpenTok.DeleteArchive(archive.ID)
}
