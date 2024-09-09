package router

import (
	"github.com/labstack/echo/v4"
	"github.com/niko-cb/uct/internal/di"
	"github.com/niko-cb/uct/internal/infrastructure/web/config"
)

// invoice is a function to create a new Resource struct for the invoice API
func invoice() *Resource {
	var invoiceHandler = di.InitializeInvoiceHandler(&config.Cfg)

	return &Resource{
		Resource: "invoices",
		Endpoints: []*Endpoint{
			{
				Method: echo.POST, SuffixPath: "", HandlerFunc: invoiceHandler.CreateInvoice,
			},
			{
				Method: echo.GET, SuffixPath: "", HandlerFunc: invoiceHandler.GetInvoicesByDateRange,
			},
		},
	}
}
