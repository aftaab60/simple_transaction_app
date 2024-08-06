//go:build wireinject
// +build wireinject

package internal

import (
	"context"
	"github.com/google/wire"
	"simple_transaction_app/internal/config"
	"simple_transaction_app/internal/db_manager"
	"simple_transaction_app/repositories"
	"simple_transaction_app/routes"
	"simple_transaction_app/services"
)

var configSet = wire.NewSet(
	config.GetConfig,
	wire.Bind(new(db_manager.IConfig), new(*config.Config)),
)

var dbSet = wire.NewSet(
	db_manager.InitPgsqlConnection,
	//db_manager.InitMysqlConnection, just one line change in dependency if we want to run our app on a different database
	wire.Bind(new(db_manager.IDB), new(*db_manager.DB)),
	wire.Bind(new(db_manager.IDbTxBeginner), new(*db_manager.DB)),
)

var repositorySet = wire.NewSet(
	repositories.NewAccountsRepository,
	repositories.NewTransactionsRepository,
)

var servicesSet = wire.NewSet(
	services.NewAccountsService,
	services.NewTransactionsService,
)

var routesSet = wire.NewSet(
	routes.NewAccountsRoutes,
	routes.NewTransactionsRoutes,
)

type App struct {
	AccountsRoutes     routes.IAccountsRoutes
	TransactionsRoutes routes.ITransactionsRoutes
}

func GetApp(ctx context.Context) (*App, func(), error) {
	wire.Build(
		configSet,
		dbSet,
		routesSet,
		repositorySet,
		servicesSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil, nil
}
