package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type databaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type Database struct {
	config               databaseConfig
	positionTableName    string
	orderTableName       string
	logTableName         string
	database             *sql.DB
	addPositionStatement *sql.Stmt
	addOrderStatement    *sql.Stmt
	addLogStatement      *sql.Stmt
}

func NewDB() (db *Database) {
	return &Database{
		databaseConfig{
			Host:     "192.168.1.250",
			Port:     "3306",
			Database: "test_trademan",
			User:     "dennis",
			Password: "c0d3mysql",
		},
		"`POSITION`",
		"`ORDER`",
		"`LOG`",
		nil,
		nil,
		nil,
		nil,
	}
}

func (db *Database) Connect() (err error) {
	db.database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", db.config.User, db.config.Password, db.config.Host, db.config.Database))
	if err != nil {
		return err
	}

	err = db.PrepareAddPositionStatement()
	if err != nil {
		return err
	}
	err = db.PrepareAddOrderStatement()
	if err != nil {
		return err
	}
	err = db.PrepareAddLogStatement()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) PrepareAddPositionStatement() (err error) {
	db.addPositionStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET Symbol=?, Status=?, Risk=?, `Size`=?, EntryPrice=?, HardStopLoss=?, Notes=?, TradingViewPlan=?, RewardRiskRatio=?, Profit=?", db.positionTableName))
	return err
}

func (db *Database) AddPosition(p Position) (TradeID int64, err error) {
	result, err := db.addPositionStatement.Exec(p.Symbol, p.Side.String(), p.Risk, p.Size, p.EntryPrice, p.HardStopLoss, p.Notes, p.TradingViewPlan, p.RewardRiskRatio, p.Profit)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *Database) GetPositions() (p []Position, err error) {
	rows, err := db.database.Query("SELECT Symbol, Status, Size FROM `POSITION`;")
	if err != nil {
		log.Printf("Error querying database %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pos Position
		if err := rows.Scan(&pos.Symbol, &pos.Status, &pos.Size); err != nil {
			log.Printf("Error querying database %v", err)
			return nil, err
		}
		p = append(p, pos)
	}
	return p, nil
}

func (db *Database) PrepareAddOrderStatement() (err error) {
	db.addOrderStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET PositionID=?, ExchangeOrderID=?, Status=?, OrderType=?, `Size`=?, TriggerPrice=?, Price=?", db.orderTableName))
	return err
}

func (db *Database) AddOrder(o Order) (OrderID int64, err error) {
	result, err := db.addOrderStatement.Exec(o.PositionID, o.ExchangeOrderID, o.Status.String(), o.OrderType.String(), o.Size, o.TriggerPrice, o.Price)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *Database) PrepareAddLogStatement() (err error) {
	db.addLogStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET PositionID=?, Source=?, Text=?", db.logTableName))
	return err
}

func (db *Database) AddLog(tradeID int64, source LogSource, text string) (err error) {
	_, err = db.addLogStatement.Exec(tradeID, source, text)
	return err
}
