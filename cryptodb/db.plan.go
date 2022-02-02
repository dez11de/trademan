package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) SavePlan(p *Plan) (err error) {
	// TODO: this should probably one transaction
	result := db.gorm.FirstOrCreate(p)
	if result.Error != nil {
		return result.Error
	}

	return result.Error
}

// TODO: this gets all plans, active and logged. Make 2 seperate functions
func (db *Database) GetPlans() (plans []Plan, err error) {
	result := db.gorm.Find(&plans)

	return plans, result.Error
}

func (db *Database) GetPlan(id uint) (plan Plan, err error) {
	result := db.gorm.Preload("Pair").Preload("Orders").Where("ID = ?", id).First(&plan)

	return plan, result.Error
}
