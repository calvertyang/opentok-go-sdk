package opentok

import (
	"fmt"
	"runtime"
)

// sdkName specifies the name of the SDK
const sdkName = "opentok-go-sdk"

// sdkVersion specifies the version of the SDK
const sdkVersion = "2.3.0"

var (
	userAgent = fmt.Sprintf("Go/%s (%s-%s) %s/%s",
		runtime.Version(),
		runtime.GOARCH,
		runtime.GOOS,
		sdkName,
		sdkVersion,
	)
)

// UserAgent returns a string containing the Go version, system architecture and OS, and the opentok-go-sdk version.
func UserAgent() string {
	return userAgent
}

// Version returns the semantic version (see http://semver.org).
func Version() string {
	return sdkVersion
}
