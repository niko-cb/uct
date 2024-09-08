//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/niko-cb/uct/internal/application/usecase"
	"github.com/niko-cb/uct/internal/controller"
	"github.com/niko-cb/uct/internal/domain/service"
	"github.com/niko-cb/uct/internal/infrastructure/gateway"
	"github.com/niko-cb/uct/internal/infrastructure/persistent/mysql"
	"github.com/niko-cb/uct/internal/infrastructure/persistent/mysql/transaction"
	"github.com/niko-cb/uct/internal/infrastructure/web/config"
	"github.com/niko-cb/uct/internal/infrastructure/web/handler"
)

func InitializeInvoiceHandler(cfg *config.Config) handler.IInvoiceHandler {
	wire.Build(
		handler.NewInvoiceHandler,
		controller.NewInvoiceController,
		usecase.NewInvoiceUsecase,
		service.NewInvoiceService,
		gateway.NewInvoiceGateway,
		mysql.NewMySQLClient,
		transaction.NewTransaction,
	)
	return &handler.InvoiceHandler{}
}
