package app

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/dep"
)

func Start(
	dbConfig fw.DBConfig,
	migrationRoot string,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
	securityPolicy fw.SecurityPolicy,
	gRpcAPIPort int,
) {
	db, err := dbConnector.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	err = dbMigrationTool.Migrate(db, migrationRoot)
	if err != nil {
		panic(err)
	}

	gRpcService, err := dep.InitGRpcService("Kgs", db, securityPolicy)
	if err != nil {
		panic(err)
	}
	gRpcService.StartAndWait(gRpcAPIPort)
}
