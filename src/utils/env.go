package util

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getString(key, fallback string) string {
	godotenv.Load()

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func getInt(key string, fallback int) int {
	godotenv.Load()

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}
