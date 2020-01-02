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

---

### Session creation, signaling, and moderation

#### Creating Sessions

To create an OpenTok Session, use the `OpenTok#CreateSession(options)` method. The `options` parameter is an struct used to specify the following:

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
You can generate a token by calling the `Session#GenerateToken(options)` method on the instance after creating it.

```go
// Generate a Token from just a session_id (fetched from a database)
token, err := ot.GenerateToken(sessionId, opentok.TokenOptions{})

// Generate a Token from a session object (returned from CreateSession)
token, err := session.GenerateToken(opentok.TokenOptions{})

// Set some options in a Token
token, err := session.GenerateToken(opentok.TokenOptions{
  Role:                   "moderator",
  ExpireTime:             time.Now().UTC().Add(7 * 24 * time.Hour).Unix(), // in one week
  Data:                   "name=Johnny",
  InitialLayoutClassList: ["focus"],
})
```

The response data is a string of token.

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
project, err := ot.ChangeProjectStatus("<PROJECT_API_KEY>", opentok.ProjectSuspended)

// Change the project status to active by project API key
project, err := ot.ChangeProjectStatus("<PROJECT_API_KEY>", opentok.ProjectActive)
```

The response data is a [project details object](#project-details-object).

#### Deleting a project

Use this method to delete a project. This prevents the use of the project API key (an any OpenTok sessions created with it).

You can also temporarily [suspend a project's API key](#changing-the-status-for-a-project-api-key).

> **Note:** You can also delete a project on your [TokBox account](https://tokbox.com/account/) page.

```go
err := ot.DeleteProject("<PROJECT_API_KEY>")
```

#### Getting information about projects

Use this method to get a project details record describing the project (or to get the records for all projects).

```go
project, err := ot.GetProjectInfo("<PROJECT_API_KEY>")
```

The response data is a [project details object](#project-details-object).

#### Generating a new project API secret

For security reasons, you may want to generate a new API secret for a project.

> **Note:** Use the new API secret for all REST API calls and with the OpenTok server-side SDKs. When you generate a new API secret, all existing [client tokens](https://tokbox.com/developer/guides/create-token/) become invalid (and they cannot be used to connect to OpenTok sessions); use the new API secret with the OpenTok server SDK to generate client tokens.

```go
project, err := ot.RefreshProjectSecret("<PROJECT_API_KEY>")
```

The response data is a [project details object](#project-details-object).

---

### Responses

#### Session Details Object

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

#### Project Details Object

```go
type Project struct {
	// The OpenTok project API key
	Id string
	// The OpenTok project API secret
	Secret string
	// Whether the project is active ("ACTIVE") or suspended ("SUSPENDED").
	Status string
	// The name, if you specified one when creating the project; or an empty string if you did not specify a name
	Name string
	// This is set to "standard" or "enterprise", and it refers to the environment a project is running on.
	Environment string
	// The time at which the project was created (a UNIX timestamp, in milliseconds)
	CreatedAt int
}
```

## Requirements

You need an OpenTok API key and API secret, which you can obtain by logging into your
[TokBox account](https://tokbox.com/account).

The OpenTok Go SDK requires Go 1.12 or higher. It may work on older versions but they are no longer tested.

## Release Notes

See the [Releases](https://github.com/calvertyang/opentok-go-sdk/releases) page for details
about each release.
