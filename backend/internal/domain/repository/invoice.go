package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/niko-cb/uct/internal/domain/entity/models"
)

// InvoiceRepository is an interface for interacting with the invoice gateway
type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, tx *sql.Tx, invoice *models.Invoice) error
	GetInvoicesByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*models.Invoice, error)
}
