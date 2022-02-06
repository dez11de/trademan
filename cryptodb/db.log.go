package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) CreateLog(l *Log) (err error) {
	result := db.Create(l)

	return result.Error
}

func (db *Database) GetLogs(PlanID uint) (logs []Log, err error) {
	result := db.Where("plan_id = ?", PlanID).Find(&logs)

	return logs, result.Error
}
