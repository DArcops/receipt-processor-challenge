package receipt

import (
	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/entity"
	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/port"
	"github.com/gin-gonic/gin"
)

// Memory storage for receipts and points.
var receiptByID map[string]entity.Receipt
var receiptPointsByID map[string]int64

func RegisterRoutes(router *gin.RouterGroup, receiptService port.ReceiptService) {
	receiptByID = make(map[string]entity.Receipt)
	receiptPointsByID = make(map[string]int64)

	controller := newReceiptController(receiptService)

	router.POST("/process", controller.createReceipt)
	router.GET("/:receipt_id/points", controller.getReceiptPoints)
}
