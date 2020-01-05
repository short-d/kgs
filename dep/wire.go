// +build wireinject

package dep

import (
	"database/sql"
	"errors"

	"github.com/google/wire"
	"github.com/short-d/app/fw"
	"github.com/short-d/app/modern/mdcli"
	"github.com/short-d/app/modern/mddb"
	"github.com/short-d/app/modern/mdemail"
	"github.com/short-d/app/modern/mdenv"
	"github.com/short-d/app/modern/mdgrpc"
	"github.com/short-d/app/modern/mdio"
	"github.com/short-d/app/modern/mdlogger"
	"github.com/short-d/app/modern/mdruntime"
	"github.com/short-d/app/modern/mdservice"
	"github.com/short-d/app/modern/mdtimer"
	"github.com/short-d/app/modern/mdtracer"
	"github.com/short-d/kgs/app/adapter/db"
	"github.com/short-d/kgs/app/adapter/rpc"
	"github.com/short-d/kgs/app/adapter/rpc/proto"
	"github.com/short-d/kgs/app/usecase"
	"github.com/short-d/kgs/app/usecase/keys"
	"github.com/short-d/kgs/app/usecase/keys/gen"
	"github.com/short-d/kgs/app/usecase/notification"
	"github.com/short-d/kgs/app/usecase/repo"
	"github.com/short-d/kgs/app/usecase/transactional"
	"github.com/short-d/kgs/dep/provider"
)

func InitCommandFactory() fw.CommandFactory {
	wire.Build(
		wire.Bind(new(fw.CommandFactory), new(mdcli.CobraFactory)),
		mdcli.NewCobraFactory,
	)
	return mdcli.CobraFactory{}
}

func InitDBConnector() fw.DBConnector {
	wire.Build(
		wire.Bind(new(fw.DBConnector), new(mddb.PostgresConnector)),
		mddb.NewPostgresConnector,
	)
	return mddb.PostgresConnector{}
}

func InitDBMigrationTool() fw.DBMigrationTool {
	wire.Build(
		wire.Bind(new(fw.DBMigrationTool), new(mddb.PostgresMigrationTool)),
		mddb.NewPostgresMigrationTool,
	)
	return mddb.PostgresMigrationTool{}
}

func InitEnvironment() fw.Environment {
	wire.Build(
		wire.Bind(new(fw.Environment), new(mdenv.GoDotEnv)),

		mdenv.NewGoDotEnv,
	)
	return mdenv.GoDotEnv{}
}

func allocatedKey() keys.AllocatedKeyRepoFactory {
	return func(tx transactional.Transaction) (repo.AllocatedKey, error) {
		sqlTx, ok := tx.(*sql.Tx)

		if !ok {
			return nil, errors.New("allocatedKeyFactory expects sql.Tx")
		}

		return db.NewAllocatedKeyTransactional(sqlTx), nil
	}
}

func availableKey() keys.AvailableKeyRepoFactory {
	return func(tx transactional.Transaction) (repo.AvailableKey, error) {
		sqlTx, ok := tx.(*sql.Tx)

		if !ok {
			return nil, errors.New("availableKeyFactory expects sql.Tx")
		}

		return db.NewAvailableKeyTransactional(sqlTx), nil
	}
}

var observabilitySet = wire.NewSet(
	wire.Bind(new(fw.Logger), new(mdlogger.Local)),
	mdlogger.NewLocal,
	mdtracer.NewLocal,
)

func InitGRpcService(
	name string,
	logLevel fw.LogLevel,
	serviceEmailAddress provider.ServiceEmailAddress,
	sqlDB *sql.DB,
	securityPolicy fw.SecurityPolicy,
	sendGridAPIKey provider.SendGridAPIKey,
	templateRootDir provider.TemplateRootDir,
	cacheSize provider.CacheSize,
) (mdservice.Service, error) {
	wire.Build(
		wire.Bind(new(fw.StdOut), new(mdio.StdOut)),
		wire.Bind(new(fw.ProgramRuntime), new(mdruntime.BuildIn)),
		wire.Bind(new(fw.Server), new(mdgrpc.GRpc)),
		wire.Bind(new(fw.GRpcAPI), new(rpc.KgsAPI)),
		wire.Bind(new(fw.EmailSender), new(mdemail.SendGrid)),
		wire.Bind(new(proto.KeyGenServer), new(rpc.KeyGenServer)),
		wire.Bind(new(notification.Notifier), new(notification.EmailNotifier)),
		wire.Bind(new(keys.Producer), new(keys.ProducerPersist)),
		wire.Bind(new(keys.Consumer), new(keys.ConsumerCached)),
		wire.Bind(new(gen.Generator), new(gen.Alphabet)),
		wire.Bind(new(transactional.Factory), new(db.FactorySQL)),

		observabilitySet,

		mdio.NewBuildInStdOut,
		mdruntime.NewBuildIn,
		mdtimer.NewTimer,
		mdgrpc.NewGRpc,
		mdservice.New,
		provider.NewSendGrid,

		rpc.NewKeyGenServer,
		rpc.NewKgsAPI,
		usecase.NewUseCase,
		provider.NewHTML,
		provider.NewEmailNotifier,
		keys.NewProducerPersist,
		provider.NewConsumer,
		keys.NewConsumerPersist,
		availableKey,
		allocatedKey,
		db.NewFactorySQL,
		gen.NewAlphabet,
		gen.NewBase62,
	)
	return mdservice.Service{}, nil
}
