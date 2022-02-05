package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) CreateOrders(o *[]Order) (err error) {
    result := db.gorm.Create(&o)

    return result.Error
}

func (db *Database) SaveOrders(o *[]Order) (err error) {
	result := db.gorm.Save(&o)
	if result.Error != nil {
		return result.Error
	}

	return result.Error
}

func (db *Database) GetOrders(PlanID uint) (orders []Order, err error) {
	result := db.gorm.Where("plan_id = ?", PlanID).Find(&orders)

	return orders, result.Error
}
