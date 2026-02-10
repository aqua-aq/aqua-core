package utils

import "strings"

func SplitOnce(s string) (string, string) {
	parts := strings.SplitN(s, ":", 2)

	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], parts[1]
}
