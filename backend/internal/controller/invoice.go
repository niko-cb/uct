package controller

import (
	"context"
	"github.com/friendsofgo/errors"
	"github.com/niko-cb/uct/internal/domain/entity/models"
	"time"

	"github.com/niko-cb/uct/internal/application/usecase"
	"github.com/niko-cb/uct/internal/domain/entity"
)

type InvoiceController struct {
	use usecase.InvoiceUsecase
}

func NewInvoiceController(use usecase.InvoiceUsecase) *InvoiceController {
	return &InvoiceController{use: use}

}

func (con *InvoiceController) CreateInvoice(ctx context.Context, invoice *entity.Invoice) error {
	// Validate the invoice with a centralized validation function
	if err := validateInvoice(invoice); err != nil {
		return errors.Wrap(err, "invoice validation failed")
	}

	// Delegate to the use case for actual creation
	if err := con.use.CreateInvoice(ctx, invoice); err != nil {
		return errors.Wrap(err, "failed to create invoice")
	}

	return nil
}

// validateInvoice is a centralized validation function for invoices
func validateInvoice(invoice *entity.Invoice) error {
	if invoice == nil {
		return errors.New("invoice is required")
	}
	if invoice.CompanyID == 0 {
		return errors.New("company_id is required")
	}
	if invoice.ClientID == 0 {
		return errors.New("client_id is required")
	}
	if invoice.IssueDate.IsZero() {
		return errors.New("issue_date is required")
	}
	if invoice.DueDate.IsZero() {
		return errors.New("due_date is required")
	}
	return nil
}

func (con *InvoiceController) GetInvoicesByDateRange(ctx context.Context, from string, to string) ([]*models.Invoice, error) {
	// Validate start_date
	fromDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, errors.Wrap(err, "invalid start_date format, expected YYYY-MM-DD")
	}

	// Validate end_date
	toDate, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil, errors.Wrap(err, "invalid end_date format, expected YYYY-MM-DD")
	}

	invoices, err := con.use.GetInvoicesByDateRange(ctx, fromDate, toDate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve invoices by date range")
	}

	return invoices, nil
}
