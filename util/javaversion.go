package util

import (
	"strconv"
	"strings"
)

func JavaVersion(mcversion string) int {
	majorversion, err := strconv.ParseInt(strings.Split(mcversion, ".")[1], 10, 32)
	if err != nil {
		return 0
	}

	if majorversion <= 12 {
		return 8
	} else if majorversion < 16 {
		return 11
	} else if majorversion == 16 {
		return 16
	} else if majorversion >= 17 {
		return 17
	}

	return 0
}
