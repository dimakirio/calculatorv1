package config

import (
	"os"
	"strconv"
)

type Config struct {
	LogLevel           string
	ComputingPower     int
	TimeAddition       int
	TimeSubtraction    int
	TimeMultiplication int
	TimeDivision       int
}

func LoadConfig() *Config {
	return &Config{
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		ComputingPower:     getEnvAsInt("COMPUTING_POWER", 1),
		TimeAddition:       getEnvAsInt("TIME_ADDITION_MS", 1000),
		TimeSubtraction:    getEnvAsInt("TIME_SUBTRACTION_MS", 1000),
		TimeMultiplication: getEnvAsInt("TIME_MULTIPLICATION_MS", 1000),
		TimeDivision:       getEnvAsInt("TIME_DIVISION_MS", 1000),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
