package app

import (
	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/dep"
	"github.com/short-d/kgs/dep/provider"
)

type Config struct {
	LogLevel            fw.LogLevel
	ServiceName         string
	ServiceEmailAddress string
	MigrationRoot       string
	GRpcAPIPort         int
	SendGridAPIKey      string
	TemplateRootDir     string
	CacheSize           int
}

// Start launches kgs service
func Start(
	config Config,
	dbConfig fw.DBConfig,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
	securityPolicy fw.SecurityPolicy,
) {
	db, err := dbConnector.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	err = dbMigrationTool.MigrateUp(db, config.MigrationRoot)
	if err != nil {
		panic(err)
	}

	gRpcService, err := dep.InitGRpcService(
		config.ServiceName,
		config.LogLevel,
		provider.ServiceEmailAddress(config.ServiceEmailAddress),
		db,
		securityPolicy,
		provider.SendGridAPIKey(config.SendGridAPIKey),
		provider.TemplateRootDir(config.TemplateRootDir),
		provider.CacheSize(config.CacheSize),
	)
	if err != nil {
		panic(err)
	}
	gRpcService.StartAndWait(config.GRpcAPIPort)
}
