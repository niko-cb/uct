package conversion

import (
	"github.com/volatiletech/sqlboiler/v4/types"
)

// ConvertToDecimal converts float64 to types.Decimal to store in the database
func ConvertToDecimal(amount float64) (types.Decimal, error) {
	var d types.Decimal
	err := d.Scan(amount)
	return d, err
}
