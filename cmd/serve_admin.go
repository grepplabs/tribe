package cmd

import (
	"context"
	"github.com/go-openapi/loads"
	"github.com/grepplabs/tribe/api/v1/server"
	"github.com/grepplabs/tribe/api/v1/server/restapi"
	"github.com/grepplabs/tribe/cmd/handlers"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/database/client"
	"github.com/grepplabs/tribe/pkg/crypto"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var serveAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Tribe Admin Server",
	Run: func(cmd *cobra.Command, args []string) {
		runtribeAdmin()
	},
}

var (
	dbConfig           = new(config.DBConfig)
	passwordBCryptCost int
)

func init() {
	serveCmd.AddCommand(serveAdminCmd)
	initServerFlags(serveAdminCmd, serverConfig)
	initCorsFlags(serveAdminCmd, server.CorsConfig)
	initLogFlags(serveAdminCmd, logConfig)

	serveAdminCmd.Flags().StringVar(&dbConfig.ConnectionURL, "db-connection-url", "postgresql://tribe:secret@localhost:5432/tribe?sslmode=disable", "data source name as connection URI e.g. postgresql://user:password@localhost:5432/dbname?sslmode=disable")
	serveAdminCmd.Flags().IntVar(&dbConfig.MaxIdleConns, "db-max-idle-conns", 2, "The maximum number of connections in the idle connection pool")
	serveAdminCmd.Flags().IntVar(&dbConfig.MaxOpenConns, "db-max-open-conns", 25, "The maximum number of open connections to the database")
	serveAdminCmd.Flags().DurationVar(&dbConfig.ConnMaxLifetime, "db-conn-max-lifetime", 0, "The maximum amount of time a connection may be reused")

	serveAdminCmd.Flags().IntVar(&passwordBCryptCost, "security-password-bcrypt-cost", crypto.DefaultBCryptCost, "BCrypt cost used for password hashing. The minimum allowable cost is 4, default is 10")

	//TODO: set following after default value is removed _ = serveAdminCmd.MarkFlagRequired("db-connection-url")
}

func runtribeAdmin() {
	logger := log.NewLogger(*logConfig)

	ctx := context.Background()
	srv, err := NewAdminServer(ctx, logger)
	if err != nil {
		logger.WithError(err).Fatalf("admin server creation failed")
	}
	defer srv.Shutdown()

	if err := srv.Serve(); err != nil {
		logger.WithError(err).Fatalf("serve failed")
	}
}

type AdminServer struct {
	*Server
	ctx    context.Context
	cancel context.CancelFunc

	logger   log.Logger
	dbClient client.Client
}

func NewAdminServer(ctx context.Context, logger log.Logger) (*AdminServer, error) {
	sCtx, cancel := context.WithCancel(ctx)

	srv := AdminServer{
		ctx:    sCtx,
		cancel: cancel,
		logger: logger,
		Server: NewServer(),
	}
	srv.configureServerFlags()
	err := srv.initDBClient()
	if err != nil {
		return nil, err
	}
	srv.SetAPI(srv.instantiateAPI())

	return &srv, nil
}

func (s *AdminServer) initDBClient() error {
	dbClient, err := client.NewSQLClient(*dbConfig)
	if err != nil {
		return errors.Wrap(err, "database initialization failure")
	}
	s.dbClient = dbClient
	return nil
}

func (s *AdminServer) instantiateAPI() *restapi.TribeAPI {
	swaggerSpec, err := loads.Embedded(server.SwaggerJSON, server.FlatSwaggerJSON)
	if err != nil {
		s.logger.WithError(err).Fatalf("loading swagger api failed")
	}

	api := restapi.NewTribeAPI(swaggerSpec)
	api.Logger = func(format string, args ...interface{}) {
		s.logger.Printf(format, args...)
	}

	// healthz
	api.HealthzGetReadyHandler = handlers.NewHealthzGetReadyHandler()
	api.HealthzGetHealthyHandler = handlers.NewHealthzGetHealthyHandler()

	// realm
	api.RealmsGetRealmHandler = handlers.NewGetRealmHandler(s.dbClient)
	api.RealmsCreateRealmHandler = handlers.NewCreateRealmHandler(s.dbClient)
	api.RealmsListRealmsHandler = handlers.NewListRealmsHandler(s.dbClient)
	api.RealmsUpdateRealmHandler = handlers.NewUpdateRealmHandler(s.dbClient)
	api.RealmsDeleteRealmHandler = handlers.NewDeleteRealmHandler(s.dbClient)

	// users
	api.UsersGetUserHandler = handlers.NewGetUserHandler(s.dbClient)
	api.UsersCreateUserHandler = handlers.NewCreateUserHandler(s.dbClient, passwordBCryptCost)
	api.UsersListUsersHandler = handlers.NewListUsersHandler(s.dbClient)
	api.UsersUpdateUserHandler = handlers.NewUpdateUserHandler(s.dbClient, passwordBCryptCost)
	api.UsersDeleteUserHandler = handlers.NewDeleteUserHandler(s.dbClient)

	return api
}
