package app

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/dep"
	"github.com/byliuyang/kgs/dep/provider"
)

type Config struct {
	ServiceName         string
	ServiceEmailAddress string
	MigrationRoot       string
	GRpcAPIPort         int
	SendGridAPIKey      string
	TemplatePattern     string
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

	err = dbMigrationTool.Migrate(db, config.MigrationRoot)
	if err != nil {
		panic(err)
	}

	gRpcService, err := dep.InitGRpcService(
		config.ServiceName,
		provider.ServiceEmailAddress(config.ServiceEmailAddress),
		db,
		securityPolicy,
		provider.SendGridAPIKey(config.SendGridAPIKey),
		provider.TemplatePattern(config.TemplatePattern),
	)
	if err != nil {
		panic(err)
	}
	gRpcService.StartAndWait(config.GRpcAPIPort)
}
