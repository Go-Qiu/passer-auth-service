/*
Package jwt is a custom inplementation of the well known JWT algorithm.  This custom implementation is to illustrate the author's understanding on hashing and the publicly known application of hashing that is commonly used in JWT-based authentication and verification protocol in many Web-based application.
*/
package jwt

import (
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/sha3"
)

var ErrEmptyToken = errors.New("[JWT]: jwt cannot be empty")
var ErrEmptyKey = errors.New("[JWT]: signing key cannot be empty")
var ErrWrongFormat = errors.New("[JWT]: wrong token format")
var ErrEmptyJWTHeader = errors.New("[JWT]: jwt header is empty.")
var ErrEmptyJWTPayload = errors.New("[JWT]: jwt payload is empty.")
var ErrEmptyJWTSignature = errors.New("[JWT]: jwt signature is empty.")

// Generate creates a JWT JSON string using the parameters passed in.
// Input parameters:
// - a is a JSON string, {"algo": "SHA3.512", "typ": "JWT"}, that indicates the hashing algorithm used for generating the authenticity code (used in later verification).  Only SHA3.512 is supported;
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
func Verify(jwt string, key string) (bool, error) {

	// exceptions handling
	if strings.TrimSpace(jwt) == "" {
		// empty jwt string
		return false, ErrEmptyToken
	}

	if strings.TrimSpace(key) == "" {
		// empty key string
		return false, ErrEmptyKey
	}

	// split the token string (by '.')
	ts := strings.Split(jwt, ".")
	if len(ts) == 0 {
		// jwt is not in the anticipated format of
		// hhhhhhhhh.pppppppppp.sssssss
		return false, ErrWrongFormat
	}

	if strings.TrimSpace(ts[0]) == "" {
		// jwt header is empty.
		return false, ErrEmptyJWTHeader
	}

	if strings.TrimSpace(ts[1]) == "" {
		// jwt payload is empty.
		return false, ErrEmptyJWTPayload
	}

	if strings.TrimSpace(ts[2]) == "" {
		// jwt signature is empty.
		return false, ErrEmptyJWTSignature
	}

	// ok.
	payloadB64 := strings.Join(ts, "")
	signatureB64 := generateSignature(payloadB64, key)

	if signatureB64 == ts[2] {
		return true, nil
	}

	return false, nil
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

// genrateSignature signs the Base64 encoded payload, payloadB64 , with the key using SHA3 512 hasing function.  It returns a Base64 encoded signature string.
func generateSignature(payloadB64 string, key string) string {

	// append the key to payloadB64
	payload := payloadB64 + key

	// generate a hash using SHA3 512 hashing function
	hasher := sha3.New512()
	hasher.Write([]byte(payload))
	h := hasher.Sum(nil)

	// convert the signature into a Base64 string.
	signs := [][]byte{}
	signs = append(signs, h)

	signsB64 := []string{}
	for sign := range b64Encode(signs) {
		signsB64 = append(signsB64, sign)
	}

	return signsB64[0]
}
