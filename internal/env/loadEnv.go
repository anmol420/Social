package env

import (
	"log"
	"os"
	"strconv"
)

func StringGetEnv(key string) string {
	value := GetEnv(key)
	if value == "" {
		log.Fatalf("Environment Variable %s is Not Set!", key)
		os.Exit(1)
	}
	return value
}

func IntegerGetEnv(key string) int {
	value := GetEnv(key)
	if value == "" {
		log.Fatalf("Environment Variable %s is Not Set!", key)
		os.Exit(1)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Environment Variable %s is Not a Valid Integer!", key)
		os.Exit(1)
	}
	return intValue
}

func BoolGetEnv(key string) bool {
	value := GetEnv(key)
	if value == "" {
		log.Fatalf("Environment Variable %s is Not Set!", key)
		os.Exit(1)
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment Variable %s is Not a Valid Boolean!", key)
		os.Exit(1)
	}
	return boolValue
}
