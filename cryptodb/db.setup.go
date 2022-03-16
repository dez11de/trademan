package cryptodb

import "log"

func (db *Database) CreateSetup(s *Setup) error {
	tx := db.Begin()

	result := db.Create(&s.Plan)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for i := range s.Orders {
		s.Orders[i].PlanID = s.Plan.ID
	}

	result = db.Create(&s.Orders)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

// 	for i := range s.Orders {
// 		s.Orders[i].LinkOrderID = fmt.Sprintf("TM-%04d-%05d-%d", s.Plan.ID, s.Orders[i].ID, s.Plan.CreatedAt.Unix())
// 	}

	result = db.Save(&s.Orders)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	var logEntry Log
	logEntry.PlanID = s.Plan.ID
	logEntry.Source = User
	logEntry.Text = "Plan created."
	result = db.Create(&logEntry)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return nil
}

func (db *Database) SaveSetup(logSource LogSource, newSetup *Setup) error {
    log.Print("Saving setup...")
	pair, err := db.GetPair(newSetup.Plan.PairID)
	if err != nil {
		return err
	}

	var oldPlan Plan
	result := db.Where("ID = ?", newSetup.Plan.ID).First(&oldPlan)
	if result.Error != nil {
		return result.Error
	}
	var oldOrders []Order
	result = db.Where("plan_id = ?", newSetup.Plan.ID).Find(&oldOrders)
	if result.Error != nil {
		return result.Error
	}

	tx := db.Begin()

    log.Print("Saving plan...")
	result = tx.Save(&newSetup.Plan)
	if result.Error != nil {
        log.Printf("And error occured saving plan: %s", result.Error)
		tx.Rollback()
		return result.Error
	}

    log.Print("Logging plan differences...")
	err = logPlanDifferences(tx, logSource, pair, oldPlan, newSetup.Plan)
    if err != nil {
        log.Printf("Error logging differences in updated plan: %s", err)
        tx.Rollback()
        return err
    }

    log.Print("Saving orders...")
	result = tx.Save(newSetup.Orders)
	if result.Error != nil {
        log.Printf("And error occured saving orders: %s", result.Error)
		tx.Rollback()
		return result.Error
	}

    log.Print("Logging order differences...")
	logOrderDifferences(tx, logSource, pair, oldOrders, newSetup.Orders)
    
    log.Print("Committing setup...")
	tx.Commit()

	return err
}
