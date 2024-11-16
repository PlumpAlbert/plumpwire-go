package models

type Client struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Balance        float64 `json:"balance,omitempty"`
	PaidToDate     float64 `json:"paid_to_date,omitempty"`
	PaymentBalance float64 `json:"payment_balance,omitempty"`
}
