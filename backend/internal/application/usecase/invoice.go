package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/niko-cb/uct/internal/domain/entity/models"
	"time"

	"github.com/niko-cb/uct/internal/infrastructure/monitor/log"

	"github.com/niko-cb/uct/internal/infrastructure/persistent/mysql/transaction"

	"github.com/niko-cb/uct/internal/domain/entity"
	"github.com/niko-cb/uct/internal/domain/repository"
	"github.com/niko-cb/uct/internal/domain/service"
)

type InvoiceUsecase interface {
	CreateInvoice(ctx context.Context, invoice *entity.Invoice) error
	GetInvoicesByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*models.Invoice, error)
}

var _ InvoiceUsecase = &invoiceUsecase{}

type invoiceUsecase struct {
	invoiceService service.InvoiceService
	transaction    repository.Transaction
}

func NewInvoiceUsecase(invoiceService service.InvoiceService, transaction repository.Transaction) InvoiceUsecase {
	return &invoiceUsecase{
		invoiceService: invoiceService,
		transaction:    transaction,
	}
}

// CreateInvoice saves invoices to the database after calculating the fee, tax, and total amount
func (u *invoiceUsecase) CreateInvoice(ctx context.Context, invoice *entity.Invoice) error {
	// calculate fee, tax, and total amount

	// 4% fee
	invoice.FeeAmount = invoice.PaymentAmount * 0.04

	// 10% tax
	invoice.TaxAmount = invoice.FeeAmount * 1.10

	// total amount = payment amount + fee amount + tax amount
	invoice.TotalAmount = invoice.PaymentAmount + invoice.FeeAmount + invoice.TaxAmount

	// Even though it's just one operation, we still want to wrap it in a transaction
	// to ensure that the operation is atomic. If the operation fails, we want to roll back
	// the entire operation to avoid any partial data being saved to the database.
	return u.transaction.DoInTx(ctx, func(ctx context.Context) error {

		// Get the transaction from the context
		tx := ctx.Value(u.transaction.(*transaction.Transaction).CtxTxKey()).(*sql.Tx)

		invoiceM, err := u.invoiceService.EntityToModel(ctx, invoice)
		if err != nil {
			log.Error(ctx, fmt.Errorf("failed to convert entity to model: %+v", err))
			return err
		}

		// Create invoices
		err = u.invoiceService.CreateInvoice(ctx, tx, invoiceM)
		if err != nil {
			log.Error(ctx, fmt.Errorf("failed to upsert invoices: %+v", err))
			return err
		}

		return nil
	})

}

// GetInvoicesByDateRange retrieves saved invoices from the database
func (u *invoiceUsecase) GetInvoicesByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*models.Invoice, error) {
	log.Info(ctx, "listing invoices")
	invoices, err := u.invoiceService.GetInvoicesByDateRange(ctx, from, to)
	if err != nil {
		log.Error(ctx, fmt.Errorf("failed to get invoices: %+v", err))
		return nil, err
	}

	return invoices, nil
}
