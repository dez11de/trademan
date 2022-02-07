package cryptodb

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AnyTime struct{} // I don't actually know if I even need this

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestShouldSavePair(t *testing.T) {
	mockdb, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer mockdb.Close()

	dialector := mysql.New(mysql.Config{Conn: mockdb, DriverName: "mysql", SkipInitializeWithVersion: true})
	gdb, err := gorm.Open(dialector, &gorm.Config{})

	db := &Database{gdb}

	mockPairs := []Pair{
		{
			ID:            33,
			Name:          "BTCUSDT",
			BaseCurrency:  "BTC",
			QuoteCurrency: "USDT",
			PriceScale:    1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            42,
			Name:          "ETHUSDT",
			BaseCurrency:  "ETH",
			QuoteCurrency: "USDT",
			PriceScale:    2,
		},
	}

	t.Run("Should INSERT Pair", func(t *testing.T) {
		for n, p := range mockPairs {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `pairs` (.+) VALUES (.+)").
				WithArgs(p.Name, p.Alias, p.Status, p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.Order.Min, p.Order.Max, p.Order.Step, AnyTime{}, AnyTime{}, p.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			err := db.CreatePair(&mockPairs[n])
			if err != nil {
				t.Errorf("received unexpected error %s", err)
			}
		}
	})
	t.Run("Should UPDATE Pair", func(t *testing.T) {
		for n, p := range mockPairs {
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE `pairs` SET .+ WHERE `id` = .+").
				WithArgs(p.Name, p.Alias, p.Status, p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.Order.Min, p.Order.Max, p.Order.Step, AnyTime{}, AnyTime{}, p.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			err := db.SavePair(&mockPairs[n])
			if err != nil {
				t.Errorf("received unexpected error %s", err)
			}
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestShouldGetPairs(t *testing.T) {
	mockdb, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer mockdb.Close()

	dialector := mysql.New(mysql.Config{Conn: mockdb, DriverName: "mysql", SkipInitializeWithVersion: true})
	gdb, err := gorm.Open(dialector, &gorm.Config{})

	db := &Database{gdb}

	mockRow := sqlmock.NewRows([]string{"ID", "Name", "Alias", "BaseCurrency", "QuoteCurrency"}).AddRow(0, "BTCUSDT", "BTCUSDT", "BTC", "USDT")

	t.Run("Should SELECT all pairs", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `pairs` ORDER BY ID ASC")).
			WillReturnRows(mockRow)

		_, err := db.GetPairs()
		if err != nil {
			t.Errorf("received unexpected error %s", err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

/*
Doesn't work, problem with % in argument
func TestShouldSearchPairs(t *testing.T) {
	mockdb, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer mockdb.Close()

	dialector := mysql.New(mysql.Config{Conn: mockdb, DriverName: "mysql", SkipInitializeWithVersion: true})
	gdb, err := gorm.Open(dialector, &gorm.Config{})

	db := &Database{gdb}

	mockSearch := "AD"
	mockRow := sqlmock.NewRows([]string{"ID", "Name", "Alias", "BaseCurrency", "QuoteCurrency"}).AddRow(15, "ADAUSDT", "ADAUSDT", "ADA", "USDT")

	t.Run("Should SEARCH for pairs", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `name` FROM `pairs` WHERE name LIKE ?")).
			WithArgs("%" + mockSearch + "%").
			WillReturnRows(mockRow)

		_, err := db.FindPairNames(mockSearch)
		if err != nil {
			t.Errorf("received unexpected error %s", err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
*/
