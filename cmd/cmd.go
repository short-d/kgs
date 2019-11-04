package cmd

import (
	"fmt"
	"os"

	"github.com/byliuyang/kgs/app"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/dep"
)

// NewRootCmd creates and initializes root command
func NewRootCmd(
	config app.Config,
	dbConfig fw.DBConfig,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
	securityPolicy fw.SecurityPolicy,
) fw.Command {
	var migrationRoot string

	cmdFactory := dep.InitCommandFactory()
	startCmd := cmdFactory.NewCommand(
		fw.CommandConfig{
			Usage: "start",
			OnExecute: func(cmd *fw.Command, args []string) {
				app.Start(
					config,
					dbConfig,
					dbConnector,
					dbMigrationTool,
					securityPolicy,
				)
			},
		},
	)
	startCmd.AddStringFlag(&migrationRoot, "migration", "app/adapter/db/migration", "migrations root directory")

	rootCmd := cmdFactory.NewCommand(
		fw.CommandConfig{
			Usage:     "kgs",
			OnExecute: func(cmd *fw.Command, args []string) {},
		},
	)
	err := rootCmd.AddSubCommand(startCmd)
	if err != nil {
		panic(err)
	}
	return rootCmd
}

// Execute runs root command
func Execute(rootCmd fw.Command) {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
