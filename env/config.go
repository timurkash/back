package env

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
)

func ReplaceEnvs(data []byte, regexp *regexp.Regexp) []byte {
	words := regexp.FindAllString(string(data), -1)
	for _, word := range words {
		b := []byte(fmt.Sprintf("${%s}", word))
		if bytes.Contains(data, b) {
			val := os.Getenv(word)
			if val != "" {
				data = bytes.ReplaceAll(data, b, []byte(val))
			}
		}
	}
	return data
}
