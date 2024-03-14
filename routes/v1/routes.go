package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/josephmakin/monerochat/handlers"
	"github.com/josephmakin/monerohub/services"
)

func SetupRoutes(router *gin.Engine) {
	donationsHandler := handlers.NewDonationsHandler(
		context.TODO(),
		services.Collections["donations"],
	)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/donations", donationsHandler.ListDonationsHandler)
		v1.GET("/donation/:id", donationsHandler.GetOneDonationHandler)
		v1.POST("/donation", donationsHandler.CreateOneDonationHandler)
		v1.POST("/transaction", donationsHandler.CallbackTransactionHandler)
	}
}
