package users

import (
	"regexp"
	"strings"
)

// Function to check if the input, v is an empty string.
// Return true if not an empty string; false if it is an empty string.
func isEmptyString(v string) bool {
	return len(strings.TrimSpace(v)) == 0
}

// Function to check if the input, v is in a valid email format.
// Return true if valid; false if not valid.
func isValidEmailFormat(v string) bool {
	pattern := regexp.MustCompile(`^[a-z0-9_]+[.-][a-z0-9]+@\w+\.([a-z0-9]{2,4}|[a-z]{2}.[a-z]{2})$`)
	return pattern.MatchString(v)
}

// Function to check if the input, v is a nil slice of strings.
// Return true when nil; false when not nil.
func isEmptyStringSlice(v []string) bool {

	if v == nil {
		// nil. empty.
		return true
	}

	// not nil.  there are some strings in the slice.
	return true
}

// Function to check if the input, v is a slice of strings
// that are valid role values.
// Valid values are "ADMIN", "AGENT", "USER"
func areValidRoles(v []string) bool {

	var status bool

	for _, element := range v {

		if isEmptyString(element) {
			// element is an empty string
			break
		}

		// ok. element is not empty.

		switch element {
		case "ADMIN", "AGENT", "USER":
			status = true
		default:
			status = false
		}
	}
	return status
}
