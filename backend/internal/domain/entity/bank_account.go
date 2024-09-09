package entity

// BankAccount represents a client's bank account information
type BankAccount struct {
	ID        int64  `json:"id"`
	ClientID  int64  `json:"client_id"`
	BankName  string `json:"bank_name"`
	Branch    string `json:"branch"`
	AccountNo string `json:"account_no"`
	Holder    string `json:"holder"`
}
