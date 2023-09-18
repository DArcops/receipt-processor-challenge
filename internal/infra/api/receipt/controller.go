package receipt

import (
	"net/http"

	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/entity"
	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/port"
	"github.com/gin-gonic/gin"
)

type receiptController struct {
	receiptService port.ReceiptService
}

func newReceiptController(receiptService port.ReceiptService) *receiptController {
	return &receiptController{
		receiptService: receiptService,
	}
}

func (rc *receiptController) createReceipt(c *gin.Context) {
	var receipt entity.Receipt

	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"The receipt is invalid": err.Error()})
		return
	}

	receiptID := rc.receiptService.CreateReceiptID(c)
	receiptByID[receiptID] = receipt

	c.JSON(http.StatusOK, gin.H{"id": receiptID})
}

func (rc *receiptController) getReceiptPoints(c *gin.Context) {
	receiptID := c.Param("receipt_id")

	receipt, ok := receiptByID[receiptID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Receipt not found for that id": receiptID})
		return
	}

	// If the points for the receipt ID are already calculated, return them.
	if receiptPointsByID[receiptID] != 0 {
		c.JSON(http.StatusOK, gin.H{"points": receiptPointsByID[receiptID]})
		return
	}

	points, err := rc.receiptService.GetReceiptPoints(c, receipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error getting receipt points": err.Error()})
		return
	}

	// Set the points for the receipt ID to avoid calculating it again.
	receiptPointsByID[receiptID] = points

	c.JSON(http.StatusOK, gin.H{"points": points})
}
