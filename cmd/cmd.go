package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/short-d/app/fw"
	"github.com/short-d/kgs/app"
	"github.com/short-d/kgs/dep"
)

// NewRootCmd creates and initializes root command
func NewRootCmd(
	config app.Config,
	dbConfig fw.DBConfig,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
	securityPolicy fw.SecurityPolicy,
	eventDispatcher fw.Dispatcher,
) fw.Command {
	var migrationRoot string

	cmdFactory := dep.InitCommandFactory()
	startCmd := cmdFactory.NewCommand(
		fw.CommandConfig{
			Usage: "start",
			OnExecute: func(cmd *fw.Command, args []string) {
				ctx, cancelFn := context.WithCancel(context.Background())

				app.Start(
					ctx,
					config,
					dbConfig,
					dbConnector,
					dbMigrationTool,
					securityPolicy,
					eventDispatcher,
				)

				listenToSystemSignals(cancelFn, func() {
					if err := eventDispatcher.Close(); err != nil {
						panic(err)
					}
				})
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

// listenToSystemSignals handles OS signals
func listenToSystemSignals(cancelFn context.CancelFunc, onInterrupt func()) {
	signalChan := make(chan os.Signal, 1)
	// listen to Interrupt
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sgn := <-signalChan
	log.Printf("Handling %s ...\n", sgn)

	cancelFn()
	onInterrupt()
}
