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
	pairTableName        string
	walletTableName      string
	positionTableName    string
	orderTableName       string
	logTableName         string
	database             *sql.DB
	addPairStatement     *sql.Stmt
	addWalletStatement   *sql.Stmt
	addPositionStatement *sql.Stmt
	addOrderStatement    *sql.Stmt
	addLogStatement      *sql.Stmt

	PairCache   map[string]Pair
	WalletCache map[string]balance
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
		"`PAIR`",
		"`WALLET`",
		"`POSITION`",
		"`ORDER`",
		"`LOG`",
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

func (db *Database) Connect() (err error) {
	db.database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db.config.User, db.config.Password, db.config.Host, db.config.Port, db.config.Database))
	if err != nil {
		return err
	}

	err = db.PrepareAddPairStatement()
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	err = db.PrepareAddWalletStatement()
	if err != nil {
		log.Printf("%v", err)
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

	db.PairCache, err = db.GetPairs()
	if err != nil {
		return err
	}
	db.WalletCache, err = db.GetRecentWallet()
	if err != nil {
		return err
	}

	return err
}
