package cryptodb

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

	pairTableName   string
	walletTableName string
	planTableName   string
	orderTableName  string
	logTableName    string
}

type Database struct {
	config             databaseConfig
	database           *sql.DB
	addPairStatement   *sql.Stmt
	addWalletStatement *sql.Stmt
	addPlanStatement   *sql.Stmt
	addOrderStatement  *sql.Stmt
	addLogStatement    *sql.Stmt

	PairCache   map[string]Pair
	WalletCache map[string]Balance
}

func NewDB() (db *Database) {
	return &Database{
		databaseConfig{
			Host:     "192.168.1.250",
			Port:     "3306",
			Database: "test_trademan",
			User:     "dennis",
			Password: "c0d3mysql",

			pairTableName:   "`PAIR`",
			walletTableName: "`WALLET`",
			planTableName:   "`PLAN`",
			orderTableName:  "`ORDER`",
			logTableName:    "`LOG`",
		},
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
	db.database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		db.config.User, db.config.Password, db.config.Host, db.config.Port, db.config.Database))
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
	err = db.PrepareAddPlanStatement()
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

func (db *Database) StorePlanAndOrders(plan Plan, orders Orders) (err error) {
	log.Printf("[db.go] Storing plan and orders")

	plan.PlanID, err = db.AddPlan(plan)

	for _, order := range orders {
		order.PlanID = plan.PlanID
		log.Printf("Storing order %v", order)
		OrderID, err := db.AddOrder(order)
		order.OrderID = OrderID
		log.Printf("Stored with PlanID %d, error was %v", OrderID, err)
	}
	return err
}
