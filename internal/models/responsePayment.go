package models

type PaymentResponse struct {
	UserId   string
	Amount   float64
	Balance  float64
	IsFailed bool
}
