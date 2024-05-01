package utils

import (
	"os"
)

func GetEnvWithDefault(key, value string) string {
	v := os.Getenv(key)
	if v == "" {
		return value
	}
	return v
}
