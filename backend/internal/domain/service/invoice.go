package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/niko-cb/uct/internal/conversion"
	"github.com/niko-cb/uct/internal/domain/repository"
	"github.com/niko-cb/uct/internal/infrastructure/monitor/log"
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
	pa, err := conversion.ConvertToDecimal(invoice.PaymentAmount)
	if err != nil {
		log.Error(ctx, fmt.Errorf(fmt.Sprintf("error converting payment amount: %v", err)))
		return nil, err
	}

	fa, err := conversion.ConvertToDecimal(invoice.FeeAmount)
	if err != nil {
		log.Error(ctx, fmt.Errorf(fmt.Sprintf("error converting fee amount: %v", err)))
		return nil, err
	}

	ta, err := conversion.ConvertToDecimal(invoice.TaxAmount)
	if err != nil {
		log.Error(ctx, fmt.Errorf(fmt.Sprintf("error converting tax amount: %v", err)))
		return nil, err
	}

	toa, err := conversion.ConvertToDecimal(invoice.TotalAmount)
	if err != nil {
		log.Error(ctx, fmt.Errorf(fmt.Sprintf("error converting total amount: %v", err)))
		return nil, err
	}

	invoiceM := &models.Invoice{
		ID:            invoice.ID,
		CompanyID:     invoice.CompanyID,
		ClientID:      invoice.ClientID,
		IssueDate:     invoice.IssueDate,
		DueDate:       invoice.DueDate,
		PaymentAmount: pa,
		FeeAmount:     fa,
		TaxAmount:     ta,
		TotalAmount:   toa,
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
