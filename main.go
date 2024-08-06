package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"simple_transaction_app/internal"
	"simple_transaction_app/internal/config"
)

func main() {
	cfg := config.GetConfig()
	log.Printf("Server starting with DB: %s", cfg.GetDatabase())

	baseCtx := context.Background()
	//setup more context fields as per app need

	app, cleanup, err := internal.GetApp(baseCtx)
	if err != nil {
		log.Fatal("error in starting app server")
	}
	if cleanup != nil {
		defer cleanup()
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // to recover from panics in execution

	setupRoutes(app, r)

	err = r.Run(fmt.Sprint(":", config.Cfg.Port))
	if err != nil {
		log.Printf("Error: %v", err)
	}

}

func setupRoutes(app *internal.App, r *gin.Engine) {
	app.AccountsRoutes.Routes(r.Group("/accounts"))
	app.TransactionsRoutes.Routes(r.Group("/transactions"))
}
