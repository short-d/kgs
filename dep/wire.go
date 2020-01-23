// +build wireinject

package dep

import (
	"database/sql"

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
	"github.com/short-d/kgs/app/usecase/dispatcher"
	"github.com/short-d/kgs/app/usecase/keys"
	"github.com/short-d/kgs/app/usecase/keys/gen"
	"github.com/short-d/kgs/app/usecase/repo"
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
	eventDispatcher fw.Dispatcher,
) (mdservice.Service, error) {
	wire.Build(
		wire.Bind(new(fw.StdOut), new(mdio.StdOut)),
		wire.Bind(new(fw.ProgramRuntime), new(mdruntime.BuildIn)),
		wire.Bind(new(fw.Server), new(mdgrpc.GRpc)),
		wire.Bind(new(fw.GRpcAPI), new(rpc.KgsAPI)),
		wire.Bind(new(fw.EmailSender), new(mdemail.SendGrid)),
		wire.Bind(new(proto.KeyGenServer), new(rpc.KeyGenServer)),
		wire.Bind(new(keys.Producer), new(keys.ProducerPersist)),
		wire.Bind(new(keys.Consumer), new(keys.ConsumerCached)),
		wire.Bind(new(gen.Generator), new(gen.Alphabet)),
		wire.Bind(new(repo.AvailableKey), new(db.AvailableKeySQL)),
		wire.Bind(new(repo.AllocatedKey), new(db.AllocatedKeySQL)),

		observabilitySet,

		// event listener subscription
		dispatcher.NewEventEmitter,
		provider.NewEmailNotifierEventListener,

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
		keys.NewProducerPersist,
		provider.NewConsumer,
		keys.NewConsumerPersist,
		db.NewAvailableKeySQL,
		db.NewAllocatedKeySQL,
		gen.NewAlphabet,
		gen.NewBase62,
	)
	return mdservice.Service{}, nil
}
