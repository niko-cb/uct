package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/niko-cb/uct/internal/controller"
	"github.com/niko-cb/uct/internal/domain/entity"
	"github.com/niko-cb/uct/internal/domain/entity/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the InvoiceUsecase
type MockInvoiceUsecase struct {
	mock.Mock
}

func (m *MockInvoiceUsecase) CreateInvoice(ctx context.Context, invoice *entity.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceUsecase) GetInvoicesByDateRange(ctx context.Context, from, to time.Time) ([]*models.Invoice, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]*models.Invoice), args.Error(1)
}

func TestCreateInvoice_ValidInvoice(t *testing.T) {
	ctx := context.Background()

	// Create mock usecase
	mockUsecase := new(MockInvoiceUsecase)

	// Create controller with the mock usecase
	c := controller.NewInvoiceController(mockUsecase)

	// Sample invoice for testing (including IssueDate and DueDate to pass validation)
	invoice := &entity.Invoice{
		ID:            1,
		CompanyID:     1,
		ClientID:      1,
		PaymentAmount: 10000.0,
		Status:        "unprocessed",
		IssueDate:     time.Now(),                          // Add IssueDate
		DueDate:       time.Now().Add(30 * 24 * time.Hour), // Add DueDate
	}

	// Set expectations for the mock
	mockUsecase.On("CreateInvoice", ctx, invoice).Return(nil)

	// Call
	err := c.CreateInvoice(ctx, invoice)

	// Assert that there are no errors and mock expectations were met
	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
}

func TestCreateInvoice_InvalidInvoice(t *testing.T) {
	ctx := context.Background()

	// Create mock usecase
	mockUsecase := new(MockInvoiceUsecase)

	// Create controller with the mock usecase
	c := controller.NewInvoiceController(mockUsecase)

	// Invalid invoice (missing CompanyID)
	invoice := &entity.Invoice{
		ID:            1,
		ClientID:      1,
		PaymentAmount: 10000.0,
		Status:        "unprocessed",
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
	}

	// Call
	err := c.CreateInvoice(ctx, invoice)

	// Adjust the expected error message to match the actual output
	assert.Error(t, err)
	assert.EqualError(t, err, "invoice validation failed: company_id is required")
}

func TestGetInvoicesByDateRange_ValidDates(t *testing.T) {
	ctx := context.Background()

	// Create mock usecase
	mockUsecase := new(MockInvoiceUsecase)

	// Create controller with the mock usecase
	c := controller.NewInvoiceController(mockUsecase)

	// Define the date range
	fromDate, _ := time.Parse("2006-01-02", "2024-01-01")
	toDate, _ := time.Parse("2006-01-02", "2024-01-31")

	// Set up mock return value
	mockUsecase.On("GetInvoicesByDateRange", ctx, fromDate, toDate).Return([]*models.Invoice{}, nil)

	// Call with valid dates
	invoices, err := c.GetInvoicesByDateRange(ctx, "2024-01-01", "2024-01-31")

	// Assert that there are no errors and the result is as expected
	assert.NoError(t, err)
	assert.NotNil(t, invoices)
	mockUsecase.AssertExpectations(t)
}

func TestGetInvoicesByDateRange_InvalidDates(t *testing.T) {
	ctx := context.Background()

	// Create mock usecase
	mockUsecase := new(MockInvoiceUsecase)

	// Create controller with the mock usecase
	c := controller.NewInvoiceController(mockUsecase)

	// Call with invalid start_date
	_, err := c.GetInvoicesByDateRange(ctx, "invalid-date", "2024-01-31")

	// Adjust the expected error message to match the detailed error output
	expectedErr := "invalid start_date format, expected YYYY-MM-DD: parsing time \"invalid-date\" as \"2006-01-02\": cannot parse \"invalid-date\" as \"2006\""
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr)
}
