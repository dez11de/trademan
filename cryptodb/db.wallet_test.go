package cryptodb

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
     "database/sql/driver"
    _ "github.com/go-sql-driver/mysql"
)

type AnyTime struct {}

func (a AnyTime) Match(v driver.Value) bool {
    _, ok := v.(time.Time)
    return ok
}

func TestShouldAddWallet(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

	mockWallet := map[string]Balance{
		"BTC": {
			Symbol: "BTC",
			Equity: decimal.NewFromFloat(1.23),
		},
		"USDT": {
			Symbol: "USDT",
			Equity: decimal.NewFromFloat(123.45),
		},
	}
	mock.ExpectExec(`INSERT INTO WALLET (.+) VALUES (.+)`).
		WithArgs(mockWallet["BTC"].Symbol,
			mockWallet["BTC"].Equity,
			mockWallet["BTC"].Available,
			mockWallet["BTC"].UsedMargin,
			mockWallet["BTC"].OrderMargin,
			mockWallet["BTC"].PositionMargin,
			mockWallet["BTC"].OCCClosingFee,
			mockWallet["BTC"].OCCFundingFee,
			mockWallet["BTC"].WalletBalance,
			mockWallet["BTC"].DailyPnL,
			mockWallet["BTC"].UnrealisedPnL,
			mockWallet["BTC"].TotalPnL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = api.AddBalance(mockWallet["BTC"])

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestShouldReturnPerformance(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

	symbol := "USDT"
	periodStart := time.Now().Add(-time.Duration(1 * 24 * time.Hour))
	mock.ExpectQuery(`
SELECT (RecentEquity - PreviousEquity) / PreviousEquity * 100 AS Performance
FROM ( SELECT ( SELECT Equity
                FROM WALLET p1
                WHERE p1.EntryTime = x.PreviousTimestamp AND Symbol = x.Symbol
            ) AS PreviousEquity,
            ( SELECT Equity
                FROM WALLET p1
                WHERE p1.EntryTime = x.RecentTimestamp AND Symbol = x.Symbol
            ) AS RecentEquity
        FROM ( SELECT Symbol, MIN(EntryTime) AS PreviousTimestamp, MAX(EntryTime) AS RecentTimestamp
                FROM WALLET
                WHERE Symbol = ".+" AND EntryTime BETWEEN '.+' AND NOW()
            ) x) x2;`).
		WithArgs(symbol, AnyTime{})

	_, err = api.GetPerformance(symbol, periodStart)

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
