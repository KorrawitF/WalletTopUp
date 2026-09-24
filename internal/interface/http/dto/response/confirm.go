package response

type Confirm struct {
	TransactionID string  `json:"transaction_id"`
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	Balance       float64 `json:"balance"`
}
