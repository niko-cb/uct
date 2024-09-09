package entity

// User represents a user in the system
type User struct {
	ID        int64  `json:"id"`
	CompanyID int64  `json:"company_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
