package usecase_test

import (
	"context"
	"github.com/niko-cb/uct/internal/conversion"
	"testing"
	"time"

	"github.com/niko-cb/uct/internal/domain/entity"
	"github.com/niko-cb/uct/internal/domain/entity/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the usecase
type MockInvoiceUsecase struct {
	mock.Mock
}

func (m *MockInvoiceUsecase) CreateInvoice(ctx context.Context, invoice *entity.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceUsecase) GetInvoicesByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*models.Invoice, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]*models.Invoice), args.Error(1)
}

// TestCreateInvoice_UseCase tests the invoice creation with mocked usecase
func TestCreateInvoice_UseCase(t *testing.T) {
	ctx := context.Background()

	// Mock usecase
	mockInvoiceUsecase := new(MockInvoiceUsecase)

	// Invoice entity that we will use for testing
	invoiceEntity := &entity.Invoice{
		ID:            1,
		CompanyID:     1,
		ClientID:      1,
		PaymentAmount: 10000.0,
		FeeAmount:     400.0,
		TaxAmount:     40.0,
		TotalAmount:   10440.0,
		Status:        "unprocessed",
	}

	// Setup mock expectations
	mockInvoiceUsecase.On("CreateInvoice", mock.Anything, invoiceEntity).Return(nil)

	// Call
	err := mockInvoiceUsecase.CreateInvoice(ctx, invoiceEntity)

	// Assertions
	assert.NoError(t, err)
	mockInvoiceUsecase.AssertExpectations(t)
}

// TestGetInvoicesByDateRange_UseCase tests the retrieval of invoices by date range using mocked usecase
func TestGetInvoicesByDateRange_UseCase(t *testing.T) {
	ctx := context.Background()

	// Mock usecase
	mockInvoiceUsecase := new(MockInvoiceUsecase)

	pa, _ := conversion.ConvertToDecimal(10000.0)
	fa, _ := conversion.ConvertToDecimal(400.0)
	ta, _ := conversion.ConvertToDecimal(40.0)
	toa, _ := conversion.ConvertToDecimal(10440.0)

	// Expected return from GetInvoicesByDateRange
	expectedInvoices := []*models.Invoice{
		{
			ID:            1,
			CompanyID:     1,
			ClientID:      1,
			PaymentAmount: pa,
			FeeAmount:     fa,
			TaxAmount:     ta,
			TotalAmount:   toa,
			Status:        "unprocessed",
		},
	}

	// Setup mock expectations
	mockInvoiceUsecase.On("GetInvoicesByDateRange", mock.Anything, mock.Anything, mock.Anything).Return(expectedInvoices, nil)

	// Call
	invoices, err := mockInvoiceUsecase.GetInvoicesByDateRange(ctx, time.Now().Add(-30*24*time.Hour), time.Now())

	// Assert no error and correct result
	assert.NoError(t, err)
	assert.Equal(t, expectedInvoices, invoices)
	mockInvoiceUsecase.AssertExpectations(t)
}
