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

### Creating Sessions

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

### Generating Tokens

Once a Session is created, you can start generating Tokens for clients to use when connecting to it.
You can generate a token by calling the `Session#GenerateToken(options)` method on the instance after creating it.

```go
// Generate a Token from just a session_id (fetched from a database)
token, err := opentok.GenerateToken(sessionId, opentok.TokenOptions{})

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

## Requirements

You need an OpenTok API key and API secret, which you can obtain by logging into your
[TokBox account](https://tokbox.com/account).

The OpenTok Go SDK requires Go 1.12 or higher. It may work on older versions but they are no longer tested.

## Release Notes

See the [Releases](https://github.com/calvertyang/opentok-go-sdk/releases) page for details
about each release.
