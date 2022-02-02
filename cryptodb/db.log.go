package cryptodb

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) SaveLog(l *Log) (err error) {
	result := db.gorm.Save(l)

	return result.Error
}

func (db *Database) GetLogs(PlanID uint) (logs []Log, err error) {
	result := db.gorm.Where("PlanID = ?", PlanID).Find(&logs)

	return logs, result.Error
}
