package main

import "os"

func main() {
}

func getEnv(varName string, defaultVal string) string {
	val := os.Getenv(varName)

	if val == "" {
		return defaultVal
	}

	return val
}
