package api

import (
	receiptapi "github.com/darcops/receipt-proccessor-challenge/internal/infra/api/receipt"
	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/port"
	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/service/receipt"
	"github.com/gin-gonic/gin"
)

var (
	receiptService port.ReceiptService
)

func registerAppRoutes(server *gin.Engine) {
	receiptService = receipt.NewReceiptService()

	apiV1 := server.Group("/api/v1")

	receiptRoutes := apiV1.Group("/receipts")

	receiptapi.RegisterRoutes(receiptRoutes, receiptService)
}
