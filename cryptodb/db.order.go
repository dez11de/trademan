package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) CreateOrders(o *[]Order) (err error) {
    result := db.gorm.Create(&o)

    return result.Error
}

func (db *Database) GetOrders(PlanID uint) (orders []Order, err error) {
	result := db.gorm.Where("PlanID = ?", PlanID).Find(&orders)

	return orders, result.Error
}
