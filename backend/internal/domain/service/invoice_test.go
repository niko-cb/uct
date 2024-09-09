package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/niko-cb/uct/internal/conversion"
	"github.com/niko-cb/uct/internal/domain/entity"
	"github.com/niko-cb/uct/internal/domain/entity/models"
	"github.com/niko-cb/uct/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the Repo
type MockInvoiceRepository struct {
	mock.Mock
}

func (m *MockInvoiceRepository) CreateInvoice(ctx context.Context, tx *sql.Tx, invoice *models.Invoice) error {
	args := m.Called(ctx, tx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) GetInvoicesByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*models.Invoice, error) {
	args := m.Called(ctx, from, to)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Invoice), args.Error(1)
}

// Test for EntityToModel
func TestEntityToModel(t *testing.T) {
	ctx := context.Background()

	// Create a sample entity.Invoice
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

	// Create the service instance
	invoiceService := service.NewInvoiceService(nil)

	// Call
	invoiceModel, err := invoiceService.EntityToModel(ctx, invoiceEntity)

	// Assertions
	assert.NoError(t, err)

	// Convert types.Decimal to float64 for comparison
	pa, _ := invoiceModel.PaymentAmount.Float64()
	fa, _ := invoiceModel.FeeAmount.Float64()
	ta, _ := invoiceModel.TaxAmount.Float64()
	toa, _ := invoiceModel.TotalAmount.Float64()

	// Compare float64 values
	assert.Equal(t, invoiceEntity.PaymentAmount, pa)
	assert.Equal(t, invoiceEntity.FeeAmount, fa)
	assert.Equal(t, invoiceEntity.TaxAmount, ta)
	assert.Equal(t, invoiceEntity.TotalAmount, toa)
	assert.Equal(t, invoiceEntity.Status, invoiceModel.Status)
}

// Test for CreateInvoice with Success
func TestCreateInvoice_Success(t *testing.T) {
	ctx := context.Background()

	// Mock repository
	mockRepo := new(MockInvoiceRepository)

	// Create the service instance
	invoiceService := service.NewInvoiceService(mockRepo)

	// Sample invoice for testing
	invoice := &entity.Invoice{
		ID:            1,
		CompanyID:     1,
		ClientID:      1,
		PaymentAmount: 10000.0,
		Status:        "unprocessed",
	}

	// Set mock expectation
	mockRepo.On("CreateInvoice", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	pa, _ := conversion.ConvertToDecimal(invoice.PaymentAmount)

	// Call
	err := invoiceService.CreateInvoice(ctx, nil, &models.Invoice{
		CompanyID:     invoice.CompanyID,
		ClientID:      invoice.ClientID,
		PaymentAmount: pa,
	})

	// Assert there is no error
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test for CreateInvoice with Error
func TestCreateInvoice_Error(t *testing.T) {
	ctx := context.Background()

	// Mock repository
	mockRepo := new(MockInvoiceRepository)

	// Create the service instance
	invoiceService := service.NewInvoiceService(mockRepo)

	// Sample invoice for testing
	invoice := &entity.Invoice{
		ID:            1,
		CompanyID:     1,
		ClientID:      1,
		PaymentAmount: 10000.0,
		Status:        "unprocessed",
	}

	// Set mock expectation to return an error
	mockRepo.On("CreateInvoice", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

	pa, _ := conversion.ConvertToDecimal(10000.0)

	// Call
	err := invoiceService.CreateInvoice(ctx, nil, &models.Invoice{
		CompanyID:     invoice.CompanyID,
		ClientID:      invoice.ClientID,
		PaymentAmount: pa,
	})

	// Assert that an error is returned
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockRepo.AssertExpectations(t)
}

// Test for GetInvoicesByDateRange with Success
func TestGetInvoicesByDateRange_Success(t *testing.T) {
	ctx := context.Background()

	// Mock repository
	mockRepo := new(MockInvoiceRepository)

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

	// Set mock expectation
	mockRepo.On("GetInvoicesByDateRange", mock.Anything, mock.Anything, mock.Anything).Return(expectedInvoices, nil)

	// Create the service instance
	invoiceService := service.NewInvoiceService(mockRepo)

	// Call
	result, err := invoiceService.GetInvoicesByDateRange(ctx, time.Now().Add(-30*24*time.Hour), time.Now())

	// Assert no error and correct result
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedInvoices, result)
	mockRepo.AssertExpectations(t)
}

// Test for GetInvoicesByDateRange with Error
func TestGetInvoicesByDateRange_Error(t *testing.T) {
	ctx := context.Background()

	// Mock repository
	mockRepo := new(MockInvoiceRepository)

	// Expected error
	expectedError := errors.New("database error")

	// Setup mock to return nil and an error
	mockRepo.On("GetInvoicesByDateRange", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	// Create the service instance
	invoiceService := service.NewInvoiceService(mockRepo)

	// Call
	result, err := invoiceService.GetInvoicesByDateRange(ctx, time.Now().Add(-30*24*time.Hour), time.Now())

	// Assert that an error is returned
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	assert.Nil(t, result) // Assert that the result is nil
	mockRepo.AssertExpectations(t)
}
