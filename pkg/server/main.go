package server

import (
	_ "database/sql"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	grpcserver "github.com/kolya59/easy_normalization/pkg/transport/grpc/server"
	pubsubserver "github.com/kolya59/easy_normalization/pkg/transport/mq/server"
	restserver "github.com/kolya59/easy_normalization/pkg/transport/rest/server"
	wsserver "github.com/kolya59/easy_normalization/pkg/transport/ws/server"
)

var opts struct {
	DBURL        string `long:"db_url" env:"DATABASE_URL" description:"DB URL" required:"true"`
	ProjectID    string `long:"projectID" env:"PROJECT_ID" required:"true" default:"trrp-virus"`
	Topic        string `long:"topic" env:"TOPIC" required:"true" default:"Students"`
	Subscription string `long:"sub" env:"SUBSCRIPTION" required:"true" default:"Students-sub"`
	RESTPort     string `long:"rest_port" env:"PORT" description:"Server port" required:"true"`
	WSPort       string `long:"ws_port" env:"WS_PORT" description:"Server port" required:"true"`
	GRPCPort     string `long:"grpc_port" env:"GRPC_PORT" description:"Server port" required:"true"`
	LogLevel     string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
}

func Start(done chan interface{}) {
	// Log initialization
	zerolog.MessageFieldName = "MESSAGE"
	zerolog.LevelFieldName = "LEVEL"
	zerolog.ErrorFieldName = "ERROR"
	zerolog.TimestampFieldName = "TIME"
	zerolog.CallerFieldName = "CALLER"
	log.Logger = log.Output(os.Stderr).With().Str("PROGRAM", "easy-normalization").Caller().Logger()

	// Parse flags
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		log.Panic().Msgf("Could not parse flags: %v", err)
	}

	level, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// соединение с БД
	log.Debug().Msg("Try to connect to database")
	if err = postgresdriver.InitDatabaseConnection(opts.DBURL); err != nil {
		log.Fatal().Msgf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err = postgresdriver.CloseConnection(); err != nil {
			log.Fatal().Msgf("Could not close db connection: %v", err)
		}
	}()

	// Start servers
	go restserver.StartServer("", opts.RESTPort, done)
	log.Info().Msg("Started REST server")
	go wsserver.StartServer("", opts.WSPort, done)
	log.Info().Msg("Started WS server")
	go pubsubserver.StartServer(opts.ProjectID, opts.Topic, opts.Subscription, done)
	log.Info().Msg("Started RabbitMQ server")
	
}
