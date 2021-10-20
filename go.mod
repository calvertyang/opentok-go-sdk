module github.com/calvertyang/opentok-go-sdk/v2

go 1.13

require (
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.7.0
)

// backward incompatible
retract (
	v2.1.1
	v2.1.0
)
