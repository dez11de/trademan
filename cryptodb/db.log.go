package cryptodb

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *api) AddLog(planID int64, source LogSource, text string) (err error) {
    
    addStmt, err := db.database.Prepare("INSERT INTO `LOG` (PlanID, Source, Text) VALUES (?, ?, ?)")
    if err != nil {
        log.Printf("error preparing Statement %v", err)
        return err 
    }
    _, err = addStmt.Exec(planID, source, text)

	return err
}
