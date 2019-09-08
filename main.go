package main

import (
	"os"
	"strconv"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/cmd"
	"github.com/byliuyang/kgs/dep"
)

func main() {
	host := getEnv("DB_HOST", "localhost")
	portStr := getEnv("DB_PORT", "5432")
	port := mustInt(portStr)
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "kgs")

	dbConfig := fw.DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbName,
	}
	dbConnector := dep.InjectDBConnector()
	dbMigrationTool := dep.InjectDBMigrationTool()

	rootCmd := cmd.NewRootCmd(dbConfig, dbConnector, dbMigrationTool)
	cmd.Execute(rootCmd)
}

func getEnv(varName string, defaultVal string) string {
	val := os.Getenv(varName)

	if val == "" {
		return defaultVal
	}

	return val
}

func mustInt(numStr string) int {
	num, err := strconv.Atoi(numStr)

	if err != nil {
		panic(err)
	}

	return num
}
