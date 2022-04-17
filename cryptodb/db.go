package cryptodb

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Username     string `flag:"user||user configured on the MySQL server instance" env:"TRADEMAN_MYSQL_USERNAME"`
	Password     string `flag:"pass||password configured on the MySQL server instance" env:"TRADEMAN_MYSQL_PASSWORD"`
	Host         string `flag:"host|127.0.0.1|the host that runs the MySQL server instance" env:"TRADEMAN_MYSQL_HOST" default:"127.0.0.1"`
	Port         string `flag:"port|3306|the port the MySQL server is listening on" env:"TRADEMAN_MYSQL_PORT" default:"3306"`
	Database     string `flag:"database|trademan|name of the database created on the MySQL server instance" env:"TRADEMAN_MYSQL_DATABASE" default:"trademan"`
	TruncTables  bool   `flag:"trunc|false|truncate most tables so you begin with an empty database. Will not truncate the pairs table. Mostly for development."`
	CreateTables bool   `flag:"create_tables|false|DANGEROUS (re)create all database tables DANGEROUS\nall existing tables will be dropped"`
}

type Database struct {
	*gorm.DB
}

func makeDSN(c DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func Connect(c DatabaseConfig) (db *Database, err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			Colorful: false,
		},
	)

	dsn := makeDSN(c)
	g, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}

	db = &Database{g}

	return db, err
}

func (db *Database) TruncTables() {
	db.Exec("TRUNCATE TABLE `assessments`")
	db.Exec("TRUNCATE TABLE `logs`")
	db.Exec("TRUNCATE TABLE `orders`")
	db.Exec("TRUNCATE TABLE `plans`")
}

func (db *Database) CreateTables() {
	// TODO: handle errors
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
	db.Migrator().DropTable(Review{})
	db.Migrator().CreateTable(Review{})
}
