package request

type Verify struct {
	UserId uint    `json:"user_id"`
	Amount float64 `json:"amount"`
	Method string  `json:"payment_method"`
}
