package opentok

import (
	"log"
	"testing"
)

const key = "<your api key here>"
const secret = "<your api secret here>"

func TestJwtToken(t *testing.T) {
	ot := New(key, secret)

	_, err := ot.jwtToken(projectToken)
	if err != nil {
		log.Fatal(err)
	}

	_, err = ot.jwtToken(accountToken)
	if err != nil {
		log.Fatal(err)
	}
}
