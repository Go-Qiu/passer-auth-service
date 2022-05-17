/*
Package jwt is a custom inplementation of the well known JWT algorithm.  This custom implementation is to illustrate the author's understanding on hashing and the publicly known application of hashing that is commonly used in JWT-based authentication and verification protocol in many Web-based application.
*/
package jwt

import (
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Generate creates a JWT JSON string using the parameters passed in.
// Input parameters:
// - a is a JSON string, {"algo": "SHA256"}, that indicates the hashing algorithm used for generating the authenticity code (used in later verification).  Only SHA256 is support now;
// - b is a JSON string that contains the payload.  The attributes supported are -
//   * "username" (string), unique User Account Id;
//   * "expiresOn" (string), a valie date time string;
//   * "roles" ([]string), Roles assigned to the User;

func Generate(header string, payload string, key string) string {

	// hash the inputs, a and b with the key passed in.
	inputs := [][]byte{}

	inputs = append(inputs, []byte(header))
	inputs = append(inputs, []byte(payload))

	t := []string{}

	for s := range b64Encode(inputs) {
		t = append(t, s)
	}

	// combine all the string in ths slice
	ts := strings.Join(t, "") + key

	hasher := sha3.New512()
	hasher.Write([]byte(ts))
	h := hasher.Sum(nil)

	signs := [][]byte{}
	signs = append(signs, h)

	signsB64 := []string{}
	for sign := range b64Encode(signs) {
		signsB64 = append(signsB64, sign)
	}

	token := strings.Join(t, ".") + "." + signsB64[0]

	return token
}

// Verify uses the passed in jwt and key to execute a check on the integrity of the jwt.
// Input parameters:
// - jwt is a JSON string
// - key is the secret key used by the service to generate the jwt
// Returns:
// an error when the integrity check fails.
// nil when the integrity check is successful.
func Verify(jwt string, key string) error {

	return nil
}

// b64Encode execute base64 encoding of each element passed into it.
// Input parameters:
// - input ([][]byte) contains all the individual element ([]byte) that needs to be encoded into base64;
// Returns :
// -
func b64Encode(input [][]byte) chan string {

	ch := make(chan string)

	go func(bs [][]byte) {

		for _, element := range bs {

			b64Outcome := base64.StdEncoding.EncodeToString(element)
			ch <- b64Outcome
		}
		close(ch)
	}(input)

	return ch
}
