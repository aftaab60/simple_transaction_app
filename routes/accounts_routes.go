package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"log"
	"net/http"
	"simple_transaction_app/models"
	"simple_transaction_app/services"
	"strconv"
)

type IAccountsRoutes interface {
	Routes(r *gin.RouterGroup)
	GetAccountByAccountID(c *gin.Context)
	CreateAccount(c *gin.Context)
}

type AccountsRoutes struct {
	AccountsService services.IAccountsService
}

var NewAccountsRoutes = wire.NewSet(
	wire.Struct(new(AccountsRoutes), "*"),
	wire.Bind(new(IAccountsRoutes), new(*AccountsRoutes)))

func (ar *AccountsRoutes) Routes(r *gin.RouterGroup) {
	accountsGroup := r.Group("/")
	//add auth middleware as this is protected API. Example below to use Auth middleware
	//accountsGroup.Use(ur.AuthMiddleware.Handle())

	accountsGroup.GET("/:account_id", ar.GetAccountByAccountID)
	accountsGroup.POST("/", ar.CreateAccount)
}

func (ar *AccountsRoutes) GetAccountByAccountID(c *gin.Context) {
	ctx := c.Request.Context()
	if id := c.Param("account_id"); id != "" {
		accountId, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, "invalid account id, must be integer number")
			return
		}
		accountDto, err := ar.AccountsService.GetAccountByID(ctx, accountId)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "error getting account details")
			return
		}
		c.JSON(http.StatusOK, accountDto)
		return
	}
	c.JSON(http.StatusBadRequest, "missing account id")
}

func (ar *AccountsRoutes) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()
	var accountCreate models.AccountCreate
	err := c.BindJSON(&accountCreate)
	if err != nil {
		log.Println("Error: incorrect request body to create account")
		c.JSON(http.StatusBadRequest, "incorrect request body")
		return
	}

	// validate input request body
	if err := accountCreate.Validate(); err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	account, err := ar.AccountsService.CreateAccount(ctx, accountCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, account)
}
