package models

type Client struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Balance        float64 `json:"balance,omitempty"`
	PaidToDate     float64 `json:"paid_to_date,omitempty"`
	PaymentBalance float64 `json:"payment_balance,omitempty"`
}

type InvoiceStatus string

const (
	Draft     InvoiceStatus = "1"
	Sent      InvoiceStatus = "2"
	Partial   InvoiceStatus = "3"
	Paid      InvoiceStatus = "4"
	Cancelled InvoiceStatus = "5"
	Reversed  InvoiceStatus = "6"
	Overdue   InvoiceStatus = "-1"
	Unpaid    InvoiceStatus = "-2"
)

type Invoice struct {
	ID       string        `json:"id,omitempty"`
	ClientID string        `json:"client_id,omitempty"`
	Amount   float64       `json:"amount,omitempty"`
	Balance  float64       `json:"balance,omitempty"`
	StatusID InvoiceStatus `json:"status_id,omitempty"`
	Date     string        `json:"date,omitempty"`
	DueDate  string        `json:"due_date,omitempty"`
}
