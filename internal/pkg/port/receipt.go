package port

import (
	"context"

	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/entity"
)

// ReceiptService is the interface that wraps the basic methods for the receipt service.
type ReceiptService interface {
	CreateReceiptID(ctx context.Context) string
	GetReceiptPoints(ctx context.Context, receipt entity.Receipt) (int64, error)
}
