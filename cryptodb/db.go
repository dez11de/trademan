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
	gorm *gorm.DB
}

func makeDSN(c databaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, "gorm_"+c.Database)
}

func Connect() (db *Database, err error) {
	// TODO: read this from env/file/cli
	dbCfg := databaseConfig{
		Host:     "192.168.1.250",
		Port:     "3306",
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

    db = &Database{}

	dsn := makeDSN(dbCfg)
    db.gorm, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
    return db, err
}

func (db *Database) RecreateTables() (err error) {
    db.gorm.Debug().Migrator().DropTable(&Pair{})
	db.gorm.Debug().Migrator().CreateTable(&Pair{})
	db.gorm.Debug().Migrator().DropTable(Plan{})
	db.gorm.Debug().Migrator().CreateTable(Plan{})
	db.gorm.Debug().Migrator().DropTable(Order{})
	db.gorm.Debug().Migrator().CreateTable(Order{})
    /*
	db.gorm.Migrator().DropTable(Log{})
	db.gorm.Migrator().CreateTable(Log{})
	db.gorm.Migrator().DropTable(Balance{})
	db.gorm.Migrator().CreateTable(Balance{})
    */
    // TODO: handle errors
    return nil
}
