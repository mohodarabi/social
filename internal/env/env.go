package env

import (
	"os"
	"strconv"
)

func GetEnvAsSting(key string, fallback string) string {
	value, exist := os.LookupEnv(key)

	if !exist {
		return fallback
	}

	return value

}

func GetEnvAsInt(key string, fallback int) int {
	value, exist := os.LookupEnv(key)

	if !exist {
		return fallback
	}

	valAsInt, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return valAsInt

}
