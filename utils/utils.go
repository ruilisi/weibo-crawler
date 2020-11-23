package utils

import (
	"bytes"
	"fmt"
	"strings"
)

func Contains(s []string, target string) bool {
	for _, element := range s {
		if strings.EqualFold(element, target) {
			return true
		}
	}
	return false
}

func CreateKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
