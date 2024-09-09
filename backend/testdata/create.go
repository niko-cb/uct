package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

type Company struct {
	ID   int64
	Name string
}

type Client struct {
	ID   int64
	Name string
}

type Invoice struct {
	CompanyID     int64
	ClientID      int64
	IssueDate     time.Time
	DueDate       time.Time
	PaymentAmount float64
	FeeAmount     float64
	TaxAmount     float64
	TotalAmount   float64
	Status        string
}

func main() {
	// Database connection setup
	db, err := sql.Open("mysql", "user:user@tcp(localhost:3306)/utc?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Generate companies, users, clients, bank accounts, and invoices within a single transaction
	err = generateTestData(db, 20, 1000, 10000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully inserted all test records in one transaction!")
}

// generateTestData generates companies, users, clients, bank accounts, and invoices in a single transaction
func generateTestData(db *sql.DB, totalCompanies, totalUsers, totalInvoices int) error {
	tx, err := db.Begin() // Start a new transaction
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Step 1: Insert companies and keep track of their IDs
	companyIDs, err := insertCompanies(tx, totalCompanies)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert companies: %v", err)
	}

	// Step 2: Insert users, each assigned to a company
	err = insertUsers(tx, totalUsers, companyIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert users: %v", err)
	}

	// Step 3: Insert clients, each client is associated with a company
	clientIDs, err := insertClients(tx, companyIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert clients: %v", err)
	}

	// Step 4: Insert bank accounts for each client
	err = insertBankAccounts(tx, clientIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert bank accounts: %v", err)
	}

	// Step 5: Generate and insert invoices for the companies
	err = insertInvoices(tx, totalInvoices, companyIDs, clientIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert invoices: %v", err)
	}

	// Commit the transaction if everything is successful
	return tx.Commit()
}

// insertCompanies inserts a number of companies and returns their IDs
func insertCompanies(tx *sql.Tx, totalCompanies int) ([]int64, error) {
	var companyIDs []int64

	for i := 1; i <= totalCompanies; i++ {
		companyName := fmt.Sprintf("Company %d", i)
		ownerName := fmt.Sprintf("Owner %d", i)
		result, err := tx.Exec("INSERT INTO companies (name, owner_name) VALUES (?, ?)", companyName, ownerName)
		if err != nil {
			return nil, fmt.Errorf("failed to insert company %d: %v", i, err)
		}

		companyID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve last inserted company ID: %v", err)
		}

		companyIDs = append(companyIDs, companyID)
	}

	return companyIDs, nil
}

// insertUsers inserts a number of users and assigns them to companies
func insertUsers(tx *sql.Tx, totalUsers int, companyIDs []int64) error {
	usersCreated := 0

	for _, companyID := range companyIDs {
		// Randomly assign users to the company
		usersInCompany := rand.Intn(totalUsers-usersCreated) + 1
		if usersCreated+usersInCompany > totalUsers {
			usersInCompany = totalUsers - usersCreated
		}

		for j := 1; j <= usersInCompany; j++ {
			userName := fmt.Sprintf("User %d", usersCreated+j)
			userEmail := fmt.Sprintf("user%d@example.com", usersCreated+j)
			password := "password123"

			_, err := tx.Exec("INSERT INTO users (company_id, name, email, password) VALUES (?, ?, ?, ?)", companyID, userName, userEmail, password)
			if err != nil {
				return fmt.Errorf("failed to insert user %d: %v", usersCreated+j, err)
			}
		}

		usersCreated += usersInCompany

		// Stop if we already assigned all users
		if usersCreated >= totalUsers {
			break
		}
	}

	return nil
}

// insertClients inserts a number of clients and returns their IDs
// Each client is associated with a company from companyIDs
func insertClients(tx *sql.Tx, companyIDs []int64) ([]int64, error) {
	var clientIDs []int64

	for _, companyID := range companyIDs {
		clientName := fmt.Sprintf("Client for Company %d", companyID)

		// Insert client with associated company_id
		result, err := tx.Exec("INSERT INTO clients (name, company_id) VALUES (?, ?)", clientName, companyID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert client for company %d: %v", companyID, err)
		}

		clientID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve last inserted client ID: %v", err)
		}

		clientIDs = append(clientIDs, clientID)
	}

	return clientIDs, nil
}

// insertBankAccounts inserts bank accounts for each client
func insertBankAccounts(tx *sql.Tx, clientIDs []int64) error {
	for _, clientID := range clientIDs {
		bankName := "Bank of Test"
		branchName := fmt.Sprintf("Branch %d", clientID)
		accountNo := fmt.Sprintf("ACC%d", clientID)
		accountName := fmt.Sprintf("Account Holder %d", clientID)

		_, err := tx.Exec("INSERT INTO bank_accounts (client_id, bank_name, branch, account_no, holder) VALUES (?, ?, ?, ?, ?)",
			clientID, bankName, branchName, accountNo, accountName)
		if err != nil {
			return fmt.Errorf("failed to insert bank account for client %d: %v", clientID, err)
		}
	}
	return nil
}

// insertInvoices generates and inserts invoices for the companies
// insertInvoices generates and inserts invoices for the companies
func insertInvoices(tx *sql.Tx, totalInvoices int, companyIDs, clientIDs []int64) error {
	rand.Seed(time.Now().UnixNano()) // Initialize the random seed

	for i := 1; i <= totalInvoices; i++ {
		companyID := companyIDs[rand.Intn(len(companyIDs))] // Random company ID for the invoice
		clientID := clientIDs[rand.Intn(len(clientIDs))]    // Random client ID from the clients list
		paymentAmount := float64(10000 + i*10)
		feeAmount := paymentAmount * 0.04
		taxAmount := feeAmount * 1.10
		totalAmount := paymentAmount + feeAmount + taxAmount
		issueDate := time.Now()

		// Generate a random due date within the next year
		daysInYear := 365
		randomDays := rand.Intn(daysInYear) // Generate a random number between 0 and 364
		dueDate := issueDate.AddDate(0, 0, randomDays)

		_, err := tx.Exec(`
            INSERT INTO invoices (company_id, client_id, issue_date, due_date, payment_amount, fee_amount, tax_amount, total_amount, status)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			companyID, clientID, issueDate, dueDate, paymentAmount, feeAmount, taxAmount, totalAmount, "unprocessed")
		if err != nil {
			return fmt.Errorf("failed to insert invoice %d: %v", i, err)
		}
	}

	return nil
}
