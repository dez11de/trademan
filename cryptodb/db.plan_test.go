package cryptodb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
)

func TestShouldAddPlan(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

	mockPlan := Plan{
		PairID: 2,
		Status: StatusPlanned,
		Side:   SideLong,
		Risk:   decimal.NewFromFloat(1.0),
        Notes: "Mock testnote.",
        TradingViewPlan: "http://tradingview.com/mock/plan12491dfe",
        RewardRiskRatio: 8.76,
        Profit: decimal.NewFromFloat(0.0),
	}
	mock.ExpectExec("INSERT INTO 'PLAN' (.+) VALUES (.+)").
		WithArgs(mockPlan.PairID, mockPlan.Status, mockPlan.Side, mockPlan.Risk, mockPlan.Notes, mockPlan.TradingViewPlan, mockPlan.RewardRiskRatio, mockPlan.Profit).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = api.AddPlan(mockPlan)

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestShouldReturnPlans (t *testing.T) {
    db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

    mock.ExpectExec("SELECT * FROM PLAN").WillReturnResult(sqlmock.NewResult(0, 3))

    _, err = api.GetPlans()

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

func TestShouldReturnPlan (t *testing.T) {
    db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

    mockPlanID := int64(1)

    mock.ExpectExec("SELECT * FROM PLAN WHERE PlanID=.+").WithArgs(mockPlanID).WillReturnResult(sqlmock.NewResult(0, 1))

    _, err = api.GetPlan(mockPlanID)

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
