package cryptodb

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) SaveOrder(o *Order) (err error) {
    log.Printf("Saving order %v", o)
	result := db.gorm.FirstOrCreate(&o)
    log.Printf("Saved order %v", o)

	return result.Error
}

func (db *Database) GetOrders(PlanID uint) (orders []Order, err error) {
	result := db.gorm.Where("PlanID = ?", PlanID).Find(&orders)

	return orders, result.Error
}
