package helpers

import (
	"strings"
)

func IsEmpty(v string) bool {

	if len(strings.TrimSpace(v)) == 0 {
		return true
	}
	return false
}
