package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) CreateOrders(o []Order) (err error) {
	result := db.Create(&o)

	return result.Error
}

func (db *Database) SaveOrders(o []Order) (err error) {
	result := db.Save(&o)

	return result.Error
}

func (db *Database) SaveOrder(o *Order) (err error) {
	result := db.Save(&o)

	return result.Error
}

func (db *Database) GetOrders(PlanID uint) (orders []Order, err error) {
	result := db.Where("plan_id = ?", PlanID).Find(&orders)
	if result.RowsAffected == 0 {
		orders = NewOrders(0)
	}

	return orders, result.Error
}
