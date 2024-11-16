package models

import (
	"strings"
	"time"
)

type InvoiceDate time.Time

func (c *InvoiceDate) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", value) //parse time
	if err != nil {
		return err
	}
	*c = InvoiceDate(t) //set result using the pointer
	return nil
}

func (c InvoiceDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02 15:04:05") + `"`), nil
}
