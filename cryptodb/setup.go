package cryptodb

type Setup struct {
	Plan   Plan
	Orders []Order
}

func NewSetup() Setup {
    return Setup{Orders: NewOrders(uint64(0))}
}
