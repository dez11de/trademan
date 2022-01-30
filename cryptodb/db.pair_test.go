package cryptodb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldWritePair(t *testing.T) {
	// Here Write means Add if new pair, Update if existing Pair.
	// The PairID is a foreign key in other tables.
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

	mockPairs := map[string]Pair{
		"BTCUSDT": {
			PairID:        33,
			Pair:          "BTCUSDT",
			BaseCurrency:  "BTC",
			QuoteCurrency: "USDT",
			PriceScale:    1,
		},
		"ETHUSDT": {
			PairID:        42,
			Pair:          "ETHUSDT",
			BaseCurrency:  "ETH",
			QuoteCurrency: "USDT",
			PriceScale:    2,
		},
	}

	t.Run("Should INSERT Pairs", func(t *testing.T) {
		for n, p := range mockPairs {
			mock.ExpectExec("INSERT INTO PAIR (.+) VALUES (.+)").
				WithArgs(p.Pair, p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.OrderSize.Min, p.OrderSize.Max, p.OrderSize.Step).
				WillReturnResult(sqlmock.NewResult(1, 1))

			_, err := api.AddPair(mockPairs[n])
			if err != nil {
				t.Errorf("received unexpected error %s", err)
			}
		}
	})

	t.Run("Should UPDATE Pairs", func(t *testing.T) {
		mock.ExpectExec("UPDATE PAIR SET (.+) WHERE PairID=.").
			WithArgs(mockPairs["ETHUSDT"].BaseCurrency, mockPairs["ETHUSDT"].QuoteCurrency, mockPairs["ETHUSDT"].PriceScale, mockPairs["ETHUSDT"].TakerFee, mockPairs["ETHUSDT"].MakerFee, mockPairs["ETHUSDT"].Leverage.Min, mockPairs["ETHUSDT"].Leverage.Max, mockPairs["ETHUSDT"].Leverage.Step, mockPairs["ETHUSDT"].Price.Min, mockPairs["ETHUSDT"].Price.Max, mockPairs["ETHUSDT"].Price.Tick, mockPairs["ETHUSDT"].OrderSize.Min, mockPairs["ETHUSDT"].OrderSize.Max, mockPairs["ETHUSDT"].OrderSize.Step, mockPairs["ETHUSDT"].PairID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		_, err := api.UpdatePair(mockPairs["ETHUSDT"])
		if err != nil {
			t.Errorf("received unexpected error %s", err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
func TestShouldGetPairs(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

	mock.ExpectQuery("SELECT * FROM PAIR ORDER BY Pair")
	_, err = api.GetPairs()
	if err != nil {
		t.Errorf("received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
