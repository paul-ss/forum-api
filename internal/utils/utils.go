package utils

import (
	"strings"
	"time"
)

func RandomSlug() string {
	return strings.ReplaceAll(time.Now().String(), " ", "")
}

func GetCurrentTime(t time.Time) time.Time {
	nt := time.Time{}
	if t == nt {
		return time.Now()
	}

	return t
}

func DESC(d bool) string {
	if d {
		return " DESC "
	}
	return " "
}