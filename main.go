package main

import (
	"strconv"

	"github.com/byliuyang/kgs/app"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/cmd"
	"github.com/byliuyang/kgs/dep"
)

func main() {
	env := dep.InitEnvironment()
	env.AutoLoadDotEnvFile()

	serviceName := env.GetEnv("SERVICE_NAME", "Kgs")
	serviceEmailAddress := env.GetEnv("SERVICE_EMAIL", "kgs@localhost")

	host := env.GetEnv("DB_HOST", "localhost")
	port := mustInt(env.GetEnv("DB_PORT", "5432"))
	user := env.GetEnv("DB_USER", "postgres")
	password := env.GetEnv("DB_PASSWORD", "password")
	dbName := env.GetEnv("DB_NAME", "kgs")

	isEncryptionEnabled := mustBool(env.GetEnv("ENABLE_ENCRYPTION", ""))

	certFilePath := env.GetEnv("CERT_FILE_PATH", "")
	keyFilePath := env.GetEnv("KEY_FILE_PATH", "")

	gRpcAPIPort := mustInt(env.GetEnv("GRPC_API_PORT", "8080"))
	sendGridAPIKey := env.GetEnv("SEND_GRID_API_KEY", "")

	CacheSize := mustInt(env.GetEnv("CACHE_SIZE", "100"))

	config := app.Config{
		ServiceName:         serviceName,
		ServiceEmailAddress: serviceEmailAddress,
		MigrationRoot:       "app/adapter/db/migration",
		GRpcAPIPort:         gRpcAPIPort,
		SendGridAPIKey:      sendGridAPIKey,
		TemplateRootDir:     "app/adapter/template",
		CacheSize:           CacheSize,
	}

	dbConfig := fw.DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbName,
	}
	dbConnector := dep.InitDBConnector()
	dbMigrationTool := dep.InitDBMigrationTool()

	securityPolicy := fw.SecurityPolicy{
		IsEncrypted:         isEncryptionEnabled,
		CertificateFilePath: certFilePath,
		KeyFilePath:         keyFilePath,
	}

	rootCmd := cmd.NewRootCmd(
		config,
		dbConfig,
		dbConnector,
		dbMigrationTool,
		securityPolicy,
	)
	cmd.Execute(rootCmd)
}

func mustInt(numStr string) int {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		panic(err)
	}
	return num
}

func mustBool(boolStr string) bool {
	boolean, err := strconv.ParseBool(boolStr)
	if err != nil {
		panic(err)
	}
	return boolean
}
