package entity

// Company represents a company entity
type Company struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerName string `json:"owner_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}
