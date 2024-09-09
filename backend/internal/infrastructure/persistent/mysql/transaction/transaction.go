package transaction

import (
	"github.com/niko-cb/uct/internal/domain/repository"
	"github.com/niko-cb/uct/internal/infrastructure/persistent/mysql"
)

type Transaction struct {
	*mysql.MySQLClient
}

func NewTransaction(c *mysql.MySQLClient) repository.Transaction {
	return &Transaction{
		MySQLClient: c,
	}
}
