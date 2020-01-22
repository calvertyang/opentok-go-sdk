# OpenTok Go SDK

[![API Reference](https://godoc.org/github.com/calvertyang/opentok-go-sdk?status.svg)](https://godoc.org/github.com/calvertyang/opentok-go-sdk)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/calvertyang/opentok-go-sdk)](https://github.com/calvertyang/opentok-go-sdk/releases/latest)
[![Build Status](https://travis-ci.org/calvertyang/opentok-go-sdk.svg?branch=master)](https://travis-ci.org/calvertyang/opentok-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/calvertyang/opentok-go-sdk)](https://goreportcard.com/report/github.com/calvertyang/opentok-go-sdk)
[![license](https://img.shields.io/github/license/calvertyang/opentok-go-sdk.svg)](https://github.com/calvertyang/opentok-go-sdk/blob/master/LICENSE)

The OpenTok Go SDK lets you generate
[sessions](https://tokbox.com/developer/guides/create-session/) and
[tokens](https://tokbox.com/developer/guides/create-token/) for
[OpenTok](http://www.tokbox.com/) applications. It also includes methods for
working with OpenTok [archives](https://tokbox.com/developer/guides/archiving),
working with OpenTok [live streaming
broadcasts](https://tokbox.com/developer/guides/broadcast/live-streaming/),
working with OpenTok [SIP interconnect](https://tokbox.com/developer/guides/sip),
and [disconnecting clients from sessions](https://tokbox.com/developer/guides/moderation/rest/).

For usage and more information, please refer to [GoDoc](https://godoc.org/github.com/calvertyang/opentok-go-sdk).

## Installing

Use `go get` to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```
go get github.com/calvertyang/opentok-go-sdk
```

To update the SDK use `go get -u` to retrieve the latest version of the SDK.

```
go get -u github.com/calvertyang/opentok-go-sdk
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
session, err := ot.CreateSession(&opentok.SessionOptions{})

// The session will the OpenTok Media Router:
session, err := ot.CreateSession(&opentok.SessionOptions{
	MediaMode: opentok.Routed,
})

// A Session with a location hint
session, err := ot.CreateSession(&opentok.SessionOptions{
	Location: "12.34.56.78",
})

// A Session with an automatic archiving
session, err := ot.CreateSession(&opentok.SessionOptions{
	ArchiveMode: opentok.AutoArchived,
	MediaMode:   opentok.Routed,
})
```

#### Generating Tokens

Once a Session is created, you can start generating Tokens for clients to use when connecting to it.
You can generate a token by calling the `OpenTok.GenerateToken(sessionID, options)` method, or by calling the `Session.GenerateToken(options)` method on the instance after creating it.

```go
// Generate a Token from just a session_id (fetched from a database)
token, err := ot.GenerateToken(sessionID, &opentok.TokenOptions{})

// Generate a Token from a session object (returned from OpenTok.CreateSession)
token, err := session.GenerateToken(&opentok.TokenOptions{})

// Set some options in a Token
token, err := session.GenerateToken(&opentok.TokenOptions{
	Role:                   opentok.Moderator,
	ExpireTime:             time.Now().UTC().Add(7 * 24 * time.Hour).Unix(), // in one week
	Data:                   "name=Johnny",
	InitialLayoutClassList: []string{"focus"},
})
```

#### Sending signals

You can send a signal to all participants in an OpenTok Session by calling the `OpenTok.SendSessionSignal(sessionID, signalData)` method.

```go
err := ot.SendSessionSignal(sessionID, &opentok.SignalData{
	Type: "foo",
	Data: "bar",
})
```

Or send a signal to a specific participant in the session by calling the `OpenTok.SendConnectionSignal(sessionID, connectionID, signalData)` method.

```go
err := ot.SendConnectionSignal(sessionID, connectionID, &opentok.SignalData{
	Type: "foo",
	Data: "bar",
})
```

This is the server-side equivalent to the signal() method in the OpenTok client SDKs. See [OpenTok signaling developer guide](https://www.tokbox.com/developer/guides/signaling/).

#### Disconnecting participants

You can disconnect participants from an OpenTok Session using the `OpenTok.ForceDisconnect(sessionID, connectionID)` method.

```go
err := ot.ForceDisconnect(sessionID, connectionID)
```

This is the server-side equivalent to the forceDisconnect() method in OpenTok.js: https://www.tokbox.com/developer/guides/moderation/js/#force_disconnect.

#### Getting stream information

You can get information on an active stream in an OpenTok session

To get information on all streams in a session, call `OpenTok.ListStreams(sessionID)`.

```go
streams, err := ot.ListStreams(sessionID)
```

To get information of a specific stream in a session, call `OpenTok.GetStream(sessionID, streamID)`.

```go
stream, err := ot.GetStream(sessionID, streamID)
```

---

### Archiving

You can only record archives of sessions that use the OpenTok Media Router (with the media mode set to routed).

You can start the recording of an OpenTok Session using the `OpenTok.StartArchive(sessionID, options)` method. Note that you can only start an Archive on a Session that has clients connected.

```go
// Start an archive
archive, err := ot.StartArchive(sessionID, &opentok.ArchiveOptions{})

// Start a named Archive
archive, err := ot.StartArchive(sessionID, &opentok.ArchiveOptions{
	Name: "Important Presentation",
})
```

You can also disable audio or video recording by setting the `HasAudio` or `HasVideo` property of the `options` parameter to `false`:

```go
// Start an audio-only Archive
archive, err := ot.StartArchive(sessionID, &opentok.ArchiveOptions{
	Name:     "Important Presentation",
	HasVideo: false,
})
```

By default, all streams are recorded to a single (composed) file. You can record the different streams in the session to individual files (instead of a single composed file) by setting the `OutputMode` option to `Opentok.Individual` when you call the `OpenTok.StartArchive()` method:

```go
// Start an Archive with individual output mode
archive, err := ot.StartArchive(sessionID, &opentok.ArchiveOptions{
	Name:       "Important Presentation",
	OutputMode: opentok.Individual,
})
```

For composed archives you can set the resolution of the archive, either SD ("640x480", the default) or HD ("1280x720").

```go
// Start an Archive with HD resolution
archive, err := ot.StartArchive(sessionID, &opentok.ArchiveOptions{
	Name:       "Important Presentation",
	Resolution: opentok.HD,
})
```

You can stop the recording of a started Archive using the `OpenTok.StopArchive(archiveID)` method. You can also do this using the `Archive.Stop()` method on the Archive instance.

```go
// Stop an Archive from an archiveID (fetched from database)
result, err := ot.StopArchive(archiveID)

// Stop an Archive from an instance (returned from Opentok.StartArchive)
result, err := archive.Stop()
```

To get an `OpenTok.Archive` instance (and all the information about it) from an archiveID, use the `OpenTok.GetArchive(archiveID)` method.

```go
archive, err := ot.GetArchive(archiveID)
```

To delete an Archive, you can call the `OpenTok.deleteArchive(archiveID)` method or the `Delete()` method of an `OpenTok.Archive` instance.

```go
// Delete an Archive from an archiveID (fetched from database)
err := ot.DeleteArchive(archiveID)

// Delete an Archive from an Archive instance, returned from the OpenTok.StartArchive(),
// OpenTok.GetArchive(), or OpenTok.ListArchives() methods
err := archive.Delete()
```

You can also get a list of all the Archives you've created (up to 1000) with your API Key. This is done using the `OpenTok.ListArchives(options)` method.

```go
// Paginate through the results via offset by 100 and count by 50
archives, err := ot.ListArchives(&opentok.ArchiveListOptions{
	Offset: 100,
	Count:  50,
})

// List archives for a specific session ID
archives, err := ot.ListArchives(&opentok.ArchiveListOptions{
	SessionID: "2_MX4xMDB-flR1-QxNzIxNX4",
})
```

Note that you can also create an automatically archived session, by passing in `OpenTok.AutoArchived` as the `ArchiveMode` option when you call the `OpenTok.CreateSession()` method (see "[Creating Sessions](#creating-sessions)" above).

For an OpenTok project, you can have OpenTok upload completed archives to an Amazon S3 bucket (or an S3-compliant storage provider) or Microsoft Azure container by calling the `OpenTok.SetArchiveStorage(options)` method.

```go
// Set an archive upload target to Amazon S3 and prevents archive files from
// being stored in the OpenTok cloud if the upload fails.
result, err := ot.SetArchiveStorage(&opentok.StorageOptions{
	Type: "s3",
	Config: opentok.AmazonS3Config{
		AccessKey: "myUsername",
		SecretKey: "myPassword",
		Bucket:    "bucketName",
	},
	Fallback: "none",
})

// Set an archive upload target to an S3-compliant storage provider
result, err := ot.SetArchiveStorage(&opentok.StorageOptions{
	Type: "s3",
	Config: opentok.AmazonS3Config{
		AccessKey: "myUsername",
		SecretKey: "myPassword",
		Bucket:    "bucketName",
		Endpoint:  "s3.example.com"
	},
})

// Set an archive upload target to Microsoft Azure
result, err := ot.SetArchiveStorage(&opentok.StorageOptions{
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

For composed archives, you can change the layout dynamically, using the `OpenTok.SetArchiveLayout(archiveID, layoutOptions)` method.

```go
archive, err := ot.SetArchiveLayout(archiveID, &opentok.Layout{
	Type: opentok.PIP,
})
```

---

### SIP interconnect

You can add an audio-only stream from an external third-party SIP gateway using the SIP Interconnect feature. This requires a SIP URI, the session ID you wish to add the audio-only stream to.

```go
sipCall, err := ot.Dial(sessionID, &opentok.DialOptions{
	SIP: &opentok.SIP{
		URI: "sip:user@sip.partner.com;transport=tls",
		From: "from@example.com",
		Headers: &opentok.SIPHeaders{
			"headerKey": "headerValue",
		},
		Auth: &opentok.SIPAuth{
			Username: "username",
			Password: "password",
		},
		Secure: true,
	},
})
```

---

### Live streaming broadcasts

_Important_: Only [routed OpenTok sessions](https://tokbox.com/developer/guides/create-session/#media-mode) support live streaming broadcasts.

To start a [live streaming broadcast](https://tokbox.com/developer/guides/broadcast/live-streaming) of an OpenTok session, call the `OpenTok.StartBroadcast(sessionID, options)` method.

```go
broadcast, err := ot.StartBroadcast(sessionID, &opentok.BroadcastOptions{
	Layout: &opentok.Layout{
		Type: opentok.VerticalPresentation,
	},
	MaxDuration: 5400,
	Outputs: &opentok.BroadcastOutputOptions{
		HLS:  &opentok.HLSConfig{},
		RTMP: []*opentok.RTMPConfig{
			&opentok.RTMPConfig{
				ID:         "foo",
				ServerURL:  "rtmps://myfooserver/myfooapp",
				StreamName: "myfoostream",
			},
			&opentok.RTMPConfig{
				ID:         "bar",
				ServerURL:  "rtmp://mybarserver/mybarapp",
				StreamName: "mybarstream",
			},
		},
	},
	Resolution: opentok.HD,
})
```

See the API reference for details on the `options` parameter.

Call the `OpenTok.StopBroadcast(broadcastID)` method to stop a live streaming broadcast.

```go
broadcast, err := ot.StopBroadcast(broadcastID)
```

You can also call the `Stop()` method of the Broadcast instance to stop a broadcast.

```go
broadcast, err := ot.StopBroadcast(broadcastID)
```

Call the `Opentok.GetBroadcast(broadcastID)` method, to get a Broadcast instance.

```go
broadcast, err := ot.GetBroadcast(broadcastID)
```

You can also get a list of all the Broadcasts you've created (up to 1000) with your API Key. This is done using the `OpenTok.ListBroadcasts(options)` method.

```go
// Paginate through the results via offset by 100 and count by 50
broadcasts, err := ot.ListBroadcasts(&opentok.BroadcastListOptions{
	Offset: 100,
	Count:  50,
})

// List broadcasts for a specific session ID
broadcasts, err := ot.ListBroadcasts(&opentok.BroadcastListOptions{
	SessionID: "2_MX4xMDB-flR1-QxNzIxNX4",
})
```

To change the broadcast layout, call the `OpenTok.SetBroadcastLayout(broadcastID, layoutOptions)` method.

```go
broadcast, err := ot.SetBroadcastLayout(broadcastID, &opentok.Layout{
	Type: opentok.PIP,
})
```

You can set the initial layout class for a client's streams by setting the layout option when you create the token for the client, using the `OpenTok.generateToken()` method. And you can change the layout classes for streams in a session by calling the `OpenTok.SetStreamClassLists(sessionID, options)` method.

```go
streams, err := ot.SetStreamClassLists(sessionID, &opentok.StreamClassOptions{
	Items: []*opentok.StreamClass{
		&opentok.StreamClass{
			ID: "8b732909-0a06-46a2-8ea8-074e64d43422",
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

#### Changing the status for a project API key

Account administrators can use this method to change a project's status. The status is either active or suspended. If a project's status is suspended, you will not be able to use the project API key (and any OpenTok sessions created with it).

You can change a project's status from active to suspended and back.

```go
// Change the project status to suspended by project API key
project, err := ot.ChangeProjectStatus("PROJECT_API_KEY", opentok.ProjectSuspended)

// Change the project status to active by project API key
project, err := ot.ChangeProjectStatus("PROJECT_API_KEY", opentok.ProjectActive)
```

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

#### Generating a new project API secret

For security reasons, you may want to generate a new API secret for a project.

> **Note:** Use the new API secret for all REST API calls and with the OpenTok server-side SDKs. When you generate a new API secret, all existing [client tokens](https://tokbox.com/developer/guides/create-token/) become invalid (and they cannot be used to connect to OpenTok sessions); use the new API secret with the OpenTok server SDK to generate client tokens.

```go
project, err := ot.RefreshProjectSecret("PROJECT_API_KEY")
```

## Requirements

You need an OpenTok API key and API secret, which you can obtain by logging into your
[TokBox account](https://tokbox.com/account).

The OpenTok Go SDK requires Go 1.12 or higher. It may work on older versions but they are no longer tested.

## Release Notes

See the [Releases](https://github.com/calvertyang/opentok-go-sdk/releases) page for details
about each release.
