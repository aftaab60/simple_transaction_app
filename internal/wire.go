//go:build wireinject
// +build wireinject

package internal

import (
	"context"
	"github.com/google/wire"
	"simple_transaction_app/internal/config"
	"simple_transaction_app/internal/db_manager"
)

var configSet = wire.NewSet(
	config.GetConfig,
	wire.Bind(new(db_manager.IConfig), new(*config.Config)),
)

var dbSet = wire.NewSet(
	db_manager.InitPgsqlConnection,
	//db_manager.InitMysqlConnection, just one line change in dependency if we want to run our app on a different database
	wire.Bind(new(db_manager.IDB), new(*db_manager.DB)))

var repositorySet = wire.NewSet()

var servicesSet = wire.NewSet()

var routesSet = wire.NewSet()

type App struct {
}

func GetApp(ctx context.Context) (*App, func(), error) {
	wire.Build(
		//configSet,
		//dbSet,
		//routesSet,
		//repositorySet,
		//servicesSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil, nil
}
