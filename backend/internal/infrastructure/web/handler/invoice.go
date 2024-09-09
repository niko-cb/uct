package handler

import (
	"context"
	"net/http"

	"github.com/niko-cb/uct/internal/domain/entity"

	"github.com/labstack/echo/v4"
	"github.com/niko-cb/uct/internal/controller"
)

type IInvoiceHandler interface {
	CreateInvoice(echo.Context) error
	GetInvoicesByDateRange(echo.Context) error
}

var _ IInvoiceHandler = &InvoiceHandler{}

type InvoiceHandler struct {
	con *controller.InvoiceController
}

func NewInvoiceHandler(con *controller.InvoiceController) IInvoiceHandler {
	return &InvoiceHandler{con: con}
}

func (h *InvoiceHandler) CreateInvoice(echo echo.Context) error {
	return withContext(echo, func(ctx context.Context) error {

		// Get JSON data from request body and bind it to invoice struct
		var invoice *entity.Invoice
		if err := echo.Bind(&invoice); err != nil {
			return echo.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		err := h.con.CreateInvoice(ctx, invoice)

		if err != nil {
			return echo.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return echo.JSON(http.StatusOK, map[string]string{"message": "success"})
	})
}

// GetInvoicesByDateRange is a handler function to get invoices by date range
func (h *InvoiceHandler) GetInvoicesByDateRange(echo echo.Context) error {
	return withContext(echo, func(ctx context.Context) error {

		from := echo.QueryParam("from")
		to := echo.QueryParam("to")

		invoices, err := h.con.GetInvoicesByDateRange(ctx, from, to)
		if err != nil {
			return echo.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return echo.JSON(http.StatusOK, invoices)
	})
}
