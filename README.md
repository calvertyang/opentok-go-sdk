# OpenTok Go SDK

The OpenTok Go SDK lets you generate
[sessions](https://tokbox.com/developer/guides/create-session/) and
[tokens](https://tokbox.com/developer/guides/create-token/) for
[OpenTok](http://www.tokbox.com/) applications. It also includes methods for
working with OpenTok [archives](https://tokbox.com/developer/guides/archiving),
working with OpenTok [live streaming
broadcasts](https://tokbox.com/developer/guides/broadcast/live-streaming/),
working with OpenTok [SIP interconnect](https://tokbox.com/developer/guides/sip),
and [disconnecting clients from sessions](https://tokbox.com/developer/guides/moderation/rest/).

## Installation

```
go get github.com/calvertyang/opentok-go-sdk
```

## Usage

### Initializing

```go
import "github.com/calvertyang/opentok-go-sdk/opentok"

ot := opentok.New(apiKey, apiSecret)
```

> For most API calls, use the API secret for the specific project in your account.
>
> This is provided on the Project page of your [TokBox Account](https://tokbox.com/account/).
>
> However, the methods of **account management** are restricted to registered administrators for the OpenTok account.
>
> To use these methods, you must use the **account** API key and secret, which is only available to the account administrator,

---

### Session creation, signaling, and moderation

#### Creating Sessions

To create an OpenTok Session, use the `OpenTok.CreateSession(options)` method. The `options` parameter is an struct used to specify the following:

* Whether the session uses the [OpenTok Media
  Router](https://tokbox.com/developer/guides/create-session/#media-mode),
  which is required for some OpenTok features (such as archiving)

* A location hint for the OpenTok server.

* Whether the session is automatically archived.

```go
// Create a session that will attempt to transmit streams directly between clients.
// If clients cannot connect, the session uses the OpenTok TURN server:
session, err := ot.CreateSession(opentok.SessionOptions{})

// The session will the OpenTok Media Router:
session, err := ot.CreateSession(opentok.SessionOptions{
	MediaMode: opentok.Routed,
})

// A Session with a location hint
session, err := ot.CreateSession(opentok.SessionOptions{
	Location: "12.34.56.78",
})

// A Session with an automatic archiving
session, err := ot.CreateSession(opentok.SessionOptions{
	ArchiveMode: opentok.AutoArchived,
	MediaMode:   opentok.Routed,
})
```

The response data is a [session details object](#session-details-object).

#### Generating Tokens

Once a Session is created, you can start generating Tokens for clients to use when connecting to it.
You can generate a token by calling the `OpenTok.GenerateToken(sessionId, options)` method, or by calling the `Session.GenerateToken(options)` method on the instance after creating it.

```go
// Generate a Token from just a session_id (fetched from a database)
token, err := ot.GenerateToken(sessionId, opentok.TokenOptions{})

// Generate a Token from a session object (returned from OpenTok.CreateSession)
token, err := session.GenerateToken(opentok.TokenOptions{})

// Set some options in a Token
token, err := session.GenerateToken(opentok.TokenOptions{
	Role:                   opentok.Moderator,
	ExpireTime:             time.Now().UTC().Add(7 * 24 * time.Hour).Unix(), // in one week
	Data:                   "name=Johnny",
	InitialLayoutClassList: []string{"focus"},
})
```

The response data is a string of token.

#### Sending signals

You can send a signal to all participants in an OpenTok Session by calling the `OpenTok.SendSessionSignal(sessionId, signalData)` method.

```go
err := ot.SendSessionSignal(sessionId, opentok.SignalData{
	Type: "foo",
	Data: "bar",
})
```

Or send a signal to a specific participant in the session by calling the `OpenTok.SendConnectionSignal(sessionId, connectionId, signalData)` method.

```go
err := ot.SendConnectionSignal(sessionId, connectionId, opentok.SignalData{
	Type: "foo",
	Data: "bar",
})
```

This is the server-side equivalent to the signal() method in the OpenTok client SDKs. See [OpenTok signaling developer guide](https://www.tokbox.com/developer/guides/signaling/).

#### Disconnecting participants

You can disconnect participants from an OpenTok Session using the `OpenTok.ForceDisconnect(sessionId, connectionId)` method.

```go
err := ot.ForceDisconnect(sessionId, connectionId)
```

This is the server-side equivalent to the forceDisconnect() method in OpenTok.js: https://www.tokbox.com/developer/guides/moderation/js/#force_disconnect.

#### Getting stream information

You can get information on an active stream in an OpenTok session

To get information on all streams in a session, call `OpenTok.ListStreams(sessionId)`.

```go
streams, err := ot.ListStreams(sessionId)
```

To get information of a specific stream in a session, call `OpenTok.GetStream(sessionId, streamId)`.

```go
stream, err := ot.GetStream(sessionId, streamId)
```

---

### Archiving

You can only record archives of sessions that use the OpenTok Media Router (with the media mode set to routed).

You can start the recording of an OpenTok Session using the `OpenTok.StartArchive(sessionId, options)` method. Note that you can only start an Archive on a Session that has clients connected.

```go
// Start an archive
archive, err := ot.StartArchive(sessionId, opentok.ArchiveOptions{})

// Start a named Archive
archive, err := ot.StartArchive(sessionId, opentok.ArchiveOptions{
	Name: "Important Presentation",
})
```

You can also disable audio or video recording by setting the `HasAudio` or `HasVideo` property of the `options` parameter to `false`:

```go
// Start an audio-only Archive
archive, err := ot.StartArchive(sessionId, opentok.ArchiveOptions{
	Name:     "Important Presentation",
	HasVideo: false,
})
```

By default, all streams are recorded to a single (composed) file. You can record the different streams in the session to individual files (instead of a single composed file) by setting the `OutputMode` option to `Opentok.Individual` when you call the `OpenTok.StartArchive()` method:

```go
// Start an Archive with individual output mode
archive, err := ot.StartArchive(sessionId, opentok.ArchiveOptions{
	Name:       "Important Presentation",
	OutputMode: opentok.Individual,
})
```

For composed archives you can set the resolution of the archive, either SD ("640x480", the default) or HD ("1280x720").

```go
// Start an Archive with HD resolution
archive, err := ot.StartArchive(sessionId, opentok.ArchiveOptions{
	Name:       "Important Presentation",
	Resolution: opentok.HD,
})
```

You can stop the recording of a started Archive using the `OpenTok.StopArchive(archiveId)` method. You can also do this using the `Archive.Stop()` method on the Archive instance.

```go
// Stop an Archive from an archiveId (fetched from database)
result, err := ot.StopArchive(archiveId)

// Stop an Archive from an instance (returned from Opentok.StartArchive)
result, err := archive.Stop()
```

To get an `OpenTok.Archive` instance (and all the information about it) from an archiveId, use the `OpenTok.GetArchive(archiveId)` method.

```go
archive, err := ot.GetArchive(archiveId)
```

To delete an Archive, you can call the `OpenTok.deleteArchive(archiveId)` method or the `Delete()` method of an `OpenTok.Archive` instance.

```go
// Delete an Archive from an archiveId (fetched from database)
err := ot.DeleteArchive(archiveId)

// Delete an Archive from an Archive instance, returned from the OpenTok.StartArchive(),
// OpenTok.GetArchive(), or OpenTok.ListArchives() methods
err := archive.Delete()
```

You can also get a list of all the Archives you've created (up to 1000) with your API Key. This is done using the `OpenTok.ListArchives(options)` method.

```go
// Paginate through the results via offset by 100 and count by 50
archives, err := ot.ListArchives(opentok.ArchiveListOptions{
	Offset: 100,
	Count: 50
})

// List archives for a specific session ID
archives, err := ot.ListArchives(opentok.ArchiveListOptions{
	SessionId: "2_MX4xMDB-flR1-QxNzIxNX4",
})
```

Note that you can also create an automatically archived session, by passing in `OpenTok.AutoArchived` as the `ArchiveMode` option when you call the `OpenTok.CreateSession()` method (see "[Creating Sessions](#creating-sessions)" above).

For an OpenTok project, you can have OpenTok upload completed archives to an Amazon S3 bucket (or an S3-compliant storage provider) or Microsoft Azure container by calling the `OpenTok.SetArchiveStorage(options)` method.

```go
// Set an archive upload target to Amazon S3 and prevents archive files from
// being stored in the OpenTok cloud if the upload fails.
result, err := ot.SetArchiveStorage(opentok.StorageOptions{
	Type: "s3",
	Config: opentok.AmazonS3Config{
		AccessKey: "myUsername",
		SecretKey: "myPassword",
		Bucket:    "bucketName",
	},
	Fallback: "none",
})

// Set an archive upload target to an S3-compliant storage provider
result, err := ot.SetArchiveStorage(opentok.StorageOptions{
	Type: "s3",
	Config: opentok.AmazonS3Config{
		AccessKey: "myUsername",
		SecretKey: "myPassword",
		Bucket:    "bucketName",
		Endpoint:  "s3.example.com"
	},
})

// Set an archive upload target to Microsoft Azure
result, err := ot.SetArchiveStorage(opentok.StorageOptions{
	Type: "azure",
	Config: opentok.AzureConfig{
		AccountName: "myAccountname",
		AccountKey:  "myAccountKey",
		Container:   "containerName",
		Domain:      "domainName" // optional
	},
})
```

Once the archive upload target is set for a project, you can also delete it by calling `OpenTok.DeleteArchiveStorage()` method.

```go
// Delete the configuration of archive storage.
err := ot.DeleteArchiveStorage()
```

For composed archives, you can change the layout dynamically, using the `OpenTok.SetArchiveLayout(archiveId, layoutOptions)` method.

```go
archive, err := ot.SetArchiveLayout(archiveId, opentok.Layout{
	Type: opentok.PIP,
})
```

---

### Live streaming broadcasts

_Important_: Only [routed OpenTok sessions](https://tokbox.com/developer/guides/create-session/#media-mode) support live streaming broadcasts.

To start a [live streaming broadcast](https://tokbox.com/developer/guides/broadcast/live-streaming) of an OpenTok session, call the `OpenTok.StartBroadcast(sessionId, options)` method.

```go
broadcast, err := ot.StartBroadcast(sessionId, opentok.BroadcastOptions{
	Layout: &opentok.Layout{
		Type: opentok.VerticalPresentation,
	},
	MaxDuration: 5400,
	Outputs: opentok.BroadcastOutputOptions{
		HLS:  opentok.HLSConfig{},
		RTMP: []*opentok.RTMPConfig{
			&opentok.RTMPConfig{
				Id:         "foo",
				ServerUrl:  "rtmps://myfooserver/myfooapp",
				StreamName: "myfoostream",
			},
			&opentok.RTMPConfig{
				Id:         "bar",
				ServerUrl:  "rtmp://mybarserver/mybarapp",
				StreamName: "mybarstream",
			},
		},
	},
	Resolution: opentok.HD,
})
```

See the API reference for details on the `options` parameter.

Call the `OpenTok.StopBroadcast(broadcastId)` method to stop a live streaming broadcast.

```go
broadcast, err := ot.StopBroadcast(broadcastId)
```

You can also call the `Stop()` method of the Broadcast instance to stop a broadcast.

```go
broadcast, err := ot.StopBroadcast(broadcastId)
```

Call the `Opentok.GetBroadcast(broadcastId)` method, to get a Broadcast instance.

```go
broadcast, err := ot.GetBroadcast(broadcastId)
```

You can also get a list of all the Broadcasts you've created (up to 1000) with your API Key. This is done using the `OpenTok.ListBroadcasts(options)` method.

```go
// Paginate through the results via offset by 100 and count by 50
broadcasts, err := ot.ListBroadcasts(opentok.BroadcastListOptions{
	Offset: 100,
	Count: 50
})

// List broadcasts for a specific session ID
broadcasts, err := ot.ListBroadcasts(opentok.BroadcastListOptions{
	SessionId: "2_MX4xMDB-flR1-QxNzIxNX4",
})
```

To change the broadcast layout, call the `OpenTok.SetBroadcastLayout(broadcastId, layoutOptions)` method.

```go
broadcast, err := ot.SetBroadcastLayout(broadcastId, opentok.Layout{
	Type: opentok.PIP,
})
```

You can set the initial layout class for a client's streams by setting the layout option when you create the token for the client, using the `OpenTok.generateToken()` method. And you can change the layout classes for streams in a session by calling the `OpenTok.SetStreamClassLists(sessionId, options)` method.

```go
streams, err := ot.SetStreamClassLists(sessionId, opentok.StreamClassOptions{
	Items: []*opentok.StreamClass{
		&opentok.StreamClass{
			Id: "8b732909-0a06-46a2-8ea8-074e64d43422",
			LayoutClassList: []string{"full"},
		},
	},
})
```

Setting the layout of a live streaming broadcast is optional. By default, live streaming broadcasts use the "best fit" layout.

---

### Account management

#### Creating a new project API key

Use this method to create an OpenTok API key and secret for a project.

> **Note:** You can also create a new project on your [TokBox account](https://tokbox.com/account/) page.

```go
// Create a new project
project, err := ot.CreateProject("")

// Create a new project with specific project name
project, err := ot.CreateProject("foo")
```

The response data is a [project details object](#project-details-object).

#### Changing the status for a project API key

Account administrators can use this method to change a project's status. The status is either active or suspended. If a project's status is suspended, you will not be able to use the project API key (and any OpenTok sessions created with it).

You can change a project's status from active to suspended and back.

```go
// Change the project status to suspended by project API key
project, err := ot.ChangeProjectStatus("PROJECT_API_KEY", opentok.ProjectSuspended)

// Change the project status to active by project API key
project, err := ot.ChangeProjectStatus("PROJECT_API_KEY", opentok.ProjectActive)
```

The response data is a [project details object](#project-details-object).

#### Deleting a project

Use this method to delete a project. This prevents the use of the project API key (an any OpenTok sessions created with it).

You can also temporarily [suspend a project's API key](#changing-the-status-for-a-project-api-key).

> **Note:** You can also delete a project on your [TokBox account](https://tokbox.com/account/) page.

```go
err := ot.DeleteProject("PROJECT_API_KEY")
```

#### Getting information about projects

Use `OpenTok.GetProject(projectAPIKey)` method to get the project details record describing the project, or use `OpenTok.ListProjects()` method to get the records for all projects.

```go
// Get all projects information
projects, err := ot.ListProjects()

// Get the project information with specific project API key
projects, err := ot.GetProject("PROJECT_API_KEY")
```

The response data is an array of [project details object](#project-details-object).

#### Generating a new project API secret

For security reasons, you may want to generate a new API secret for a project.

> **Note:** Use the new API secret for all REST API calls and with the OpenTok server-side SDKs. When you generate a new API secret, all existing [client tokens](https://tokbox.com/developer/guides/create-token/) become invalid (and they cannot be used to connect to OpenTok sessions); use the new API secret with the OpenTok server SDK to generate client tokens.

```go
project, err := ot.RefreshProjectSecret("PROJECT_API_KEY")
```

The response data is a [project details object](#project-details-object).

---

### Responses

#### Session

```go
type Session struct {
	// The session id of the project
	SessionId string
	// The OpenTok project API key
	ProjectId string
	// The creation date
	CreateDt string
	// The URL of the OpenTok media router used by the session
	MediaServerURL string
	// The instance of OpenTok
	OpenTok *OpenTok
}
```

#### Project

```go
type Project struct {
	// The OpenTok project API key
	Id string
	// The OpenTok account id
	UserId int
	// The OpenTok project API secret
	Secret string
	// Whether the project is active ("VALID", "ACTIVE") or suspended ("SUSPENDED").
	Status string
	// The OpenTok account status
	UserStatus string
	// The name, if you specified one when creating the project; or an empty string if you did not specify a name
	Name string
	// The OpenTok account email
	ContactEmail string
	// The time at which the project was created (a UNIX timestamp, in milliseconds)
	CreatedAt int
	// The time at which the project was updated (a UNIX timestamp, in milliseconds)
	UpdatedAt int
	// The environment id that project is running on
	EnvironmentId int
	// The environment name that project is running on
	EnvironmentName string
	// The environment description that project is running on
	EnvironmentDescription string
	// The OpenTok project API key
	ApiKey string
}
```

#### Archive

```go
type Archive struct {
	// The time at which the archive was created, in milliseconds since the UNIX epoch.
	CreatedAt int
	// The duration of the archive, in milliseconds.
	Duration int
	// Whether the archive has an audio track or not.
	HasAudio bool
	// Whether the archive has an video track or not.
	HasVideo bool
	// The unique archive ID.
	Id string
	// The name of the archive.
	Name *string
	// The output mode to be generated for this archive.
	OutputMode ArchiveOutputMode
	// The API key associated with the archive.
	ProjectId int
	// This string describes the reason the archive stopped or failed.
	Reason string
	// The resolution of the archive.
	Resolution Resolution
	// The session ID of the OpenTok session associated with this archive.
	SessionId string
	// The size of the MP4 file.
	Size int
	// The status of the archive.
	Status string
	// The download URL of the available MP4 file.
	Url *string
	// The instance of OpenTok
	OpenTok *OpenTok
}
```

#### Archive List

```go
type ArchiveList struct {
	// The total number of archives for the API key.
	Count int
	// An array of archive defining each archive retrieved.
	Items []*Archive
}
```

#### Archive Storage Options

```go
type StorageOptions struct {
	// Type of storage.
	Type string
	// Settings for the storage.
	Config interface{}
	// Error handling method if upload fails.
	Fallback string
}
```

#### Stream

```go
type Stream struct {
	// The unique stream ID.
	Id string
	// Either "camera" or "screen".
	VideoType string
	// The stream name.
	Name string
	// An array of the layout classes for the stream.
	LayoutClassList []string
}
```

#### Stream List

```go
type StreamList struct {
	// The total number of streams for the session.
	Count int
	// An array of stream defining each stream retrieved.
	Items []*Stream
}
```

#### Broadcast

```go
type Broadcast struct {
	// The unique ID for the broadcast.
	Id string
	// The OpenTok session ID.
	SessionId string
	// The API key associated with the broadcast..
	ProjectId int
	// The time at which the broadcast was created, in milliseconds since the UNIX epoch.
	CreatedAt int
	// The time at which the broadcast was updated, in milliseconds since the UNIX epoch.
	UpdatedAt int
	// The resolution of the broadcast.
	Resolution Resolution
	// The status of the broadcast.
	Status string
	// An object containing details about the HLS and RTMP broadcasts.
	BroadcastUrls BroadcastUrls
	// The instance of OpenTok
	OpenTok *OpenTok
}
```

#### Broadcast List

```go
type BroadcastList struct {
	// The total number of broadcasts for the session.
	Count int
	// An array of broadcast defining each broadcast retrieved.
	Items []*Broadcast
}
```

---

### Type Definition

#### Project Status

```go
type ProjectStatus string

const (
	/**
	 * Set to ACTIVE to use the project API key.
	 */
	ProjectActive ProjectStatus = "ACTIVE"
	/**
	 * Set to SUSPENDED, you will not be able to use the project API key (and any OpenTok sessions created with it).
	 */
	ProjectSuspended ProjectStatus = "SUSPENDED"
)
```

#### Archive Mode

```go
type ArchiveMode string

const (
	/**
	 * Set to always to have the session archived automatically.
	 */
	AutoArchived ArchiveMode = "always"
	/**
	 * Set to manual (the default), you can archive the session by calling the REST /archive POST method
	 */
	ManualArchived ArchiveMode = "manual"
)
```

#### Media Mode

```go
type MediaMode string

const (
	/**
	 * Set to enabled if you prefer clients to attempt to send audio-video streams directly to other clients
	 */
	Relayed MediaMode = "enabled"
	/**
	 * Set to disabled for sessions that use the OpenTok Media Router
	 */
	Routed MediaMode = "disabled"
)
```

#### Role

```go
type Role string

const (
	/**
	 * A publisher can publish streams, subscribe to streams, and signal.
	 */
	Publisher Role = "publisher"
	/**
	 * A subscriber can only subscribe to streams.
	 */
	Subscriber Role = "subscriber"
	/**
	 * In addition to the privileges granted to a publisher, in clients using
	 * the OpenTok.js library, a moderator can call the forceUnpublish() and
	 * forceDisconnect() method of the Session object.
	 */
	Moderator Role = "moderator"
)
```

#### Layout Type

```go
type LayoutType string

const (
	/**
	 * This is a tiled layout, which scales according to the number of videos.
	 */
	BestFit LayoutType = "bestFit"
	/**
	 * This is a picture-in-picture layout, where a small stream is visible over
	 * a full-size stream.
	 */
	PIP LayoutType = "pip"
	/**
	 * This is a layout with one large stream on the right edge of the output,
	 * and several smaller streams along the left edge of the output.
	 */
	VerticalPresentation LayoutType = "verticalPresentation"
	/**
	 * This is a layout with one large stream on the top edge of the output,
	 * and several smaller streams along the bottom edge of the output.
	 */
	HorizontalPresentation LayoutType = "horizontalPresentation"
	/**
	 * To use a custom layout, set the type property for the layout to "custom"
	 * and set an additional property, stylesheet, which is set to the CSS.
	 */
	Custom LayoutType = "custom"
)
```

#### Archive Output Mode

```go
type ArchiveOutputMode string

const (
	/**
	 * The archive is a single MP4 file composed of all streams.
	 */
	Composed ArchiveOutputMode = "composed"
	/**
	 * The archive is a ZIP container file with multiple individual media files
	 * for each stream, and a JSON metadata file for video synchronization.
	 */
	Individual ArchiveOutputMode = "individual"
)
```

#### Resolution

```go
type Resolution string

const (
	// The resolution of the archive.
	SD Resolution = "640x480"
	HD Resolution = "1280x720"
)
```

## Requirements

You need an OpenTok API key and API secret, which you can obtain by logging into your
[TokBox account](https://tokbox.com/account).

The OpenTok Go SDK requires Go 1.12 or higher. It may work on older versions but they are no longer tested.

## Release Notes

See the [Releases](https://github.com/calvertyang/opentok-go-sdk/releases) page for details
about each release.
