// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func GetApp(ctx context.Context) (*App, func(), error) {
	configConfig := config.GetConfig()
	db := db_manager.InitPgsqlConnection(configConfig)
	accountsRepository := &repositories.AccountsRepository{
		DB: db,
	}
	accountsService := &services.AccountsService{
		AccountsRepository: accountsRepository,
	}
	accountsRoutes := &routes.AccountsRoutes{
		AccountsService: accountsService,
	}
	transactionsRepository := &repositories.TransactionsRepository{
		DB: db,
	}
	transactionsService := &services.TransactionsService{
		TransactionsRepository: transactionsRepository,
		AccountsRepository:     accountsRepository,
		Db:                     db,
	}
	transactionsRoutes := &routes.TransactionsRoutes{
		TransactionsService: transactionsService,
	}
	app := &App{
		AccountsRoutes:     accountsRoutes,
		TransactionsRoutes: transactionsRoutes,
	}
	return app, func() {
	}, nil
}

// wire.go:

var configSet = wire.NewSet(config.GetConfig, wire.Bind(new(db_manager.IConfig), new(*config.Config)))

var dbSet = wire.NewSet(db_manager.InitPgsqlConnection, wire.Bind(new(db_manager.IDB), new(*db_manager.DB)), wire.Bind(new(db_manager.IDbTxBeginner), new(*db_manager.DB)))

var repositorySet = wire.NewSet(repositories.NewAccountsRepository, repositories.NewTransactionsRepository)

var servicesSet = wire.NewSet(services.NewAccountsService, services.NewTransactionsService)

var routesSet = wire.NewSet(routes.NewAccountsRoutes, routes.NewTransactionsRoutes)

type App struct {
	AccountsRoutes     routes.IAccountsRoutes
	TransactionsRoutes routes.ITransactionsRoutes
}
