package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/niko-cb/uct/internal/domain/repository"
	"github.com/niko-cb/uct/internal/infrastructure/monitor/log"
	"github.com/volatiletech/sqlboiler/v4/types"
	"time"

	"github.com/niko-cb/uct/internal/domain/entity"
	"github.com/niko-cb/uct/internal/domain/entity/models"
)

type InvoiceService interface {
	EntityToModel(ctx context.Context, invoice *entity.Invoice) (*models.Invoice, error)
	CreateInvoice(ctx context.Context, tx *sql.Tx, invoice *models.Invoice) error
	GetInvoicesByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*models.Invoice, error)
}

type invoiceService struct {
	repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) InvoiceService {
	return &invoiceService{
		repo: repo,
	}
}

// EntityToModel converts invoice entity to an invoice model to prepare for database insert
func (s *invoiceService) EntityToModel(ctx context.Context, invoice *entity.Invoice) (*models.Invoice, error) {

	invoiceM := &models.Invoice{
		ID:            invoice.ID,
		CompanyID:     invoice.CompanyID,
		ClientID:      invoice.ClientID,
		IssueDate:     invoice.IssueDate,
		DueDate:       invoice.DueDate,
		PaymentAmount: convertToDecimal(invoice.PaymentAmount),
		FeeAmount:     convertToDecimal(invoice.FeeAmount),
		TaxAmount:     convertToDecimal(invoice.TaxAmount),
		TotalAmount:   convertToDecimal(invoice.TotalAmount),
		Status:        invoice.Status,
	}
	return invoiceM, nil
}

// CreateInvoice saves invoices to the database
func (s *invoiceService) CreateInvoice(ctx context.Context, tx *sql.Tx, invoice *models.Invoice) error {
	return s.repo.CreateInvoice(ctx, tx, invoice)
}

// GetInvoicesByDateRange retrieves invoices from the database by date range
func (s *invoiceService) GetInvoicesByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*models.Invoice, error) {
	return s.repo.GetInvoicesByDateRange(ctx, from, to)
}

// convertToDecimal converts float64 to types.Decimal to store in the database
func convertToDecimal(amount float64) types.Decimal {
	var d types.Decimal
	err := d.Scan(amount)
	if err != nil {
		log.Fatal(context.Background(), fmt.Errorf("failed to convert float64 to decimal: %+v", err))
	}

	return d
}
