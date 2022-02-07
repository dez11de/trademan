package cryptodb

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type databaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type Database struct {
	*gorm.DB
}

func makeDSN(c databaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, "gorm_"+c.Database)
}

func Connect() (db *Database, err error) {
	// TODO: read this from env/file/cli
	dbCfg := databaseConfig{
		Host: "192.168.1.250",
		Port: "3306",
		//TODO: make a switch for test/production/dev?
		Database: "test_trademan",
		User:     "dennis",
		Password: "c0d3mysql",
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			Colorful: false,
		},
	)

	dsn := makeDSN(dbCfg)
	g, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Panicf("unable to connect to database")
	}
	db = &Database{g}
	return db, err
}

func (db *Database) RecreateTables() (err error) {
	db.Migrator().DropTable(&Pair{})
	db.Migrator().CreateTable(&Pair{})
	db.Migrator().DropTable(Plan{})
	db.Migrator().CreateTable(Plan{})
	db.Migrator().DropTable(Order{})
	db.Migrator().CreateTable(Order{})
	db.Migrator().DropTable(Log{})
	db.Migrator().CreateTable(Log{})
	db.Migrator().DropTable(Balance{})
	db.Migrator().CreateTable(Balance{})
	// TODO: handle errors
	return nil
}
