package cryptodb

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) AddBalance(b Balance) (err error) {
	result := db.gorm.Save(b)

	return result.Error
}

// TODO: rewrite this into a GetBalance (symbol string) Balance, err function
func (db *Database) GetCurrentBalance(s string) (balance Balance, err error) {
	result := db.gorm.Where("Symbol = ?", s).Last(&balance)

	return balance, result.Error
}

func (db *Database) GetPerformance(s string, periodStart time.Time) (performance float64, err error) {
	var currentBalance, previousBalance Balance
	result := db.gorm.Where("Symbol = ?", s).Last(&currentBalance)
	if result.Error != nil {
		return 0.0, result.Error
	}

	result = db.gorm.Limit(1).Order("created_at").Where("Symbol = ?", s).Where("created_at <= ?", periodStart).First(&previousBalance)

	performance = previousBalance.Equity.Sub(currentBalance.Equity).Div(previousBalance.Equity).InexactFloat64()
	return performance, err
}
