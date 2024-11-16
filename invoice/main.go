package invoice

import (
	"encoding/json"
	"errors"
	"net/http"

	"plumpalbert.xyz/plumpwire/invoice/models"
)

type InvoiceManager struct {
	endpoint string

	Clients []models.Client
}

// Generate new Invoice object
func New(host string) (*InvoiceManager, error) {
	invoice := InvoiceManager{
		endpoint: host + "/api/v1",
		Clients:  []models.Client{},
	}

	err := invoice.GetClients()
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// Get list of clients
func (inv InvoiceManager) GetClients() error {
	res, err := http.Get(inv.endpoint + "/clients?status=active")
	if err != nil {
		return err
	}

	return json.NewDecoder(res.Body).Decode(&inv.Clients)
}

// Get client object by name
func (inv InvoiceManager) GetClient(client_name string) (*models.Client, error) {
	err := inv.GetClients()
	if err != nil {
		return nil, err
	}

	for _, c := range inv.Clients {
		if c.Name == client_name {
			return &c, nil
		}
	}

	return nil, errors.New("could not find client `" + client_name + "`")
}
