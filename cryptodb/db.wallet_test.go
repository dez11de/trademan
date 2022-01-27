package cryptodb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
)

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
		WithArgs( mockWallet["BTC"].Symbol,
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

	err = api.AddWallet(mockWallet["BTC"])

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
