package utils

import "math/rand"

func RandomUserAgent() string {
	userAgent := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36",
	}
	return userAgent[rand.Intn(len(userAgent))]
}
