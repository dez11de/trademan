package cryptodb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
)

func TestShouldAddOrder(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	api := NewDB()
	api.database = db

	mockOrder := Order{
		PlanID:    2,
		Status:    StatusPlanned,
		OrderType: TypeEntry,
        ExchangeOrderID: "bbtrademan-2-2",
		Size:      decimal.NewFromFloat(45.0),
		Price:     decimal.NewFromFloat(1.0263),
	}
	mock.ExpectExec("INSERT INTO 'ORDER' (.+) VALUES (.+)").
		WithArgs(mockOrder.PlanID, mockOrder.ExchangeOrderID, mockOrder.Status, mockOrder.OrderType, mockOrder.Size, mockOrder.TriggerPrice, mockOrder.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = api.AddOrder(mockOrder)

	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}
