package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) CreatePlan(p *Plan) (err error) {
	result := db.gorm.Create(&p)

	return result.Error
}

func (db *Database) SavePlan(p *Plan) (err error) {
	result := db.gorm.Save(&p)

	return result.Error
}

// TODO: this gets all plans; active and logged. Make 2 seperate functions or use a scope
// See: https://gorm.io/docs/scopes.html for ideas
func (db *Database) GetPlans() (plans []Plan, err error) {
	result := db.gorm.Find(&plans)

	return plans, result.Error
}

func (db *Database) GetPlan(id uint) (plan Plan, err error) {
	result := db.gorm.Where("ID = ?", id).First(&plan)

	return plan, result.Error
}
