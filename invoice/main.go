package invoice

import (
	"plumpalbert.xyz/plumpwire/invoice/models"
)

type Invoice struct {
	endpoint string

	Clients []models.Client
}
