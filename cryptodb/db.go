package cryptodb

import (
	"database/sql"
	"fmt"

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

type api struct {
	config             databaseConfig
	database           *sql.DB
}

func NewDB() (db *api) {
	return &api{
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
	}
}

func (db *api) Connect() (err error) {
	db.database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		db.config.User, db.config.Password, db.config.Host, db.config.Port, db.config.Database))

	return err
}
