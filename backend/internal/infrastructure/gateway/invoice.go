package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/niko-cb/uct/internal/infrastructure/monitor/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"

	"github.com/niko-cb/uct/internal/domain/entity/models"
	"github.com/niko-cb/uct/internal/domain/repository"
	"github.com/niko-cb/uct/internal/infrastructure/persistent/mysql"
)

var _ repository.InvoiceRepository = &invoiceGateway{}

type invoiceGateway struct {
	client *mysql.MySQLClient
}

func NewInvoiceGateway(client *mysql.MySQLClient) repository.InvoiceRepository {
	return &invoiceGateway{
		client: client,
	}
}

func (g *invoiceGateway) CreateInvoice(ctx context.Context, tx *sql.Tx, invoice *models.Invoice) error {
	// Connect to the database
	g.client.Connect()

	// Insert the invoice
	err := invoice.Insert(ctx, tx, boil.Infer())
	if err != nil {
		log.Error(ctx, fmt.Errorf(fmt.Sprintf("failed to insert invoice into database: %+v", err)))
		return err
	}

	return nil
}

func (g *invoiceGateway) GetInvoicesByDateRange(ctx context.Context, from, to time.Time) ([]*models.Invoice, error) {
	// Ensure the database connection is established
	g.client.Connect()

	// Retrieve invoices where the DueDate is between the provided 'from' and 'to' dates
	invoices, err := models.Invoices(
		qm.Where("due_date >= ? AND due_date <= ?", from, to),
	).All(ctx, g.client.DB)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}
