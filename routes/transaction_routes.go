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

type ITransactionsRoutes interface {
	Routes(r *gin.RouterGroup)
	GetTransactionByTransactionID(c *gin.Context)
	CreateTransaction(c *gin.Context)
}

type TransactionsRoutes struct {
	TransactionsService services.ITransactionsService
}

var NewTransactionsRoutes = wire.NewSet(
	wire.Struct(new(TransactionsRoutes), "*"),
	wire.Bind(new(ITransactionsRoutes), new(*TransactionsRoutes)))

func (tr *TransactionsRoutes) Routes(r *gin.RouterGroup) {
	transactionsGroup := r.Group("/")
	//add auth middleware as this is protected API. Example below to use Auth middleware
	//transactionsGroup.Use(ur.AuthMiddleware.Handle())

	transactionsGroup.GET("/:transaction_id", tr.GetTransactionByTransactionID)
	transactionsGroup.POST("/", tr.CreateTransaction)
}

func (tr *TransactionsRoutes) GetTransactionByTransactionID(c *gin.Context) {
	ctx := c.Request.Context()
	if id := c.Param("transaction_id"); id != "" {
		transactionId, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, "invalid transaction id, must be integer number")
			return
		}
		transactionDto, err := tr.TransactionsService.GetTransactionById(ctx, transactionId)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "error getting transaction details")
			return
		}
		c.JSON(http.StatusOK, transactionDto)
		return
	}
	c.JSON(http.StatusBadRequest, "missing transaction id")
}

func (tr *TransactionsRoutes) CreateTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	var transactionCreate models.TransactionCreate
	err := c.BindJSON(&transactionCreate)
	if err != nil {
		log.Println("Error: incorrect request body to create transaction")
		c.JSON(http.StatusBadRequest, "incorrect request body")
		return
	}

	// validate input request body
	if err := transactionCreate.Validate(); err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := tr.TransactionsService.CreateTransaction(ctx, transactionCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, transaction)
}
