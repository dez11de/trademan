package cryptodb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldAddPair(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

	mockPair := Pair{
		PairID:        1,
		Pair:          "BTCUSDT",
		BaseCurrency:  "BTC",
		QuoteCurrency: "USDT",
		PriceScale:    1,
	}
	mock.ExpectExec("INSERT INTO 'PAIR' (.+) VALUES (.+)").
		WithArgs(mockPair.Pair, mockPair.BaseCurrency, mockPair.QuoteCurrency, mockPair.PriceScale, mockPair.TakerFee, mockPair.MakerFee, mockPair.Leverage.Min, mockPair.Leverage.Max, mockPair.Leverage.Step, mockPair.Price.Min, mockPair.Price.Max, mockPair.Price.Tick, mockPair.OrderSize.Min, mockPair.OrderSize.Max, mockPair.OrderSize.Step).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = api.AddPair(mockPair)

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
