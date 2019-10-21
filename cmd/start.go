package cmd

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/dep"
)

func start(
	dbConfig fw.DBConfig,
	migrationRoot string,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
) {
	db, err := dbConnector.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	err = dbMigrationTool.Migrate(db, migrationRoot)
	if err != nil {
		panic(err)
	}

	gRpcService, err := dep.InjectGRpcService("Kgs", db)
	if err != nil {
		panic(err)
	}
	gRpcService.StartAndWait(8080)
}
