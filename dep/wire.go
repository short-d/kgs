// +build wireinject

package dep

import (
	"database/sql"

	"github.com/byliuyang/kgs/app/adapter/rpc/proto"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/app/modern/mdcli"
	"github.com/byliuyang/app/modern/mddb"
	"github.com/byliuyang/app/modern/mdemail"
	"github.com/byliuyang/app/modern/mdenv"
	"github.com/byliuyang/app/modern/mdgrpc"
	"github.com/byliuyang/app/modern/mdlogger"
	"github.com/byliuyang/app/modern/mdservice"
	"github.com/byliuyang/app/modern/mdtracer"
	"github.com/byliuyang/kgs/app/adapter/db"
	"github.com/byliuyang/kgs/app/adapter/rpc"
	"github.com/byliuyang/kgs/app/usecase/keys"
	"github.com/byliuyang/kgs/app/usecase/keys/gen"
	"github.com/byliuyang/kgs/app/usecase/notification"
	"github.com/byliuyang/kgs/app/usecase/repo"
	"github.com/byliuyang/kgs/dep/provider"
	"github.com/google/wire"
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
	mdlogger.NewLocal,
	mdtracer.NewLocal,
)

func InitGRpcService(
	name string,
	serviceEmailAddress provider.ServiceEmailAddress,
	sqlDB *sql.DB,
	securityPolicy fw.SecurityPolicy,
	sendGridAPIKey provider.SendGridAPIKey,
	templatePattern provider.TemplatePattern,
) (mdservice.Service, error) {
	wire.Build(
		wire.Bind(new(fw.Server), new(mdgrpc.GRpc)),
		wire.Bind(new(fw.GRpcAPI), new(rpc.KgsAPI)),
		wire.Bind(new(fw.EmailSender), new(mdemail.SendGrid)),
		wire.Bind(new(proto.KeyGenServer), new(rpc.KeyGenServer)),
		wire.Bind(new(notification.Notifier), new(notification.EmailNotifier)),
		wire.Bind(new(keys.Producer), new(keys.ProducerPersist)),
		wire.Bind(new(keys.Consumer), new(keys.ConsumerPersist)),
		wire.Bind(new(gen.Generator), new(gen.Alphabet)),
		wire.Bind(new(repo.AvailableKey), new(db.AvailableKeySQL)),
		wire.Bind(new(repo.AllocatedKey), new(db.AllocatedKeySQL)),

		observabilitySet,

		mdgrpc.NewGRpc,
		mdservice.New,
		provider.NewSendGrid,

		rpc.NewKeyGenServer,
		rpc.NewKgsAPI,
		provider.NewEmailNotifier,
		provider.NewTemplate,
		keys.NewProducerPersist,
		keys.NewConsumerPersist,
		db.NewAvailableKeySQL,
		db.NewAllocatedKeySQL,
		gen.NewAlphabet,
		gen.NewBase62,
	)
	return mdservice.Service{}, nil
}
