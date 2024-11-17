package invoice

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"plumpalbert.xyz/plumpwire/invoice/models"
)

type InvoiceManager struct {
	endpoint string
	c        *httpClient

	Clients []models.Client
}

// Generate new Invoice object
func New(host string, token string) (*InvoiceManager, error) {
	invoice := InvoiceManager{
		endpoint: host + "/api/v1",
		c: &httpClient{
			c:     http.Client{},
			token: token,
		},

		Clients: []models.Client{},
	}

	err := invoice.GetClients()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &invoice, nil
}

// Get list of clients
func (inv *InvoiceManager) GetClients() error {
	res, err := inv.c.Get(inv.endpoint + "/clients?status=active")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	var test struct {
		Data []models.Client `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&test)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	inv.Clients = test.Data
	slog.Debug(fmt.Sprintf("[GetClients] %#v", inv.Clients))
	return nil
}

// Get client object by name
func (inv *InvoiceManager) GetClient(client_name string) (*models.Client, error) {
	err := inv.GetClients()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for _, c := range inv.Clients {
		slog.Debug(fmt.Sprintf("[GetClient] Check: %#v == %s", c, client_name))
		if c.Name == client_name {
			slog.Debug(fmt.Sprintf("[GetClient] Found %#v", c))
			return &c, nil
		}
	}

	return nil, errors.New("could not find client `" + client_name + "`")
}

// Get list of invoices for client
func (inv *InvoiceManager) GetBills(client_name string) ([]models.Invoice, error) {
	client, err := inv.GetClient(client_name)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	res, err := inv.c.Get(inv.endpoint + "/invoices?status=active&client_status=unpaid&client_id=" + client.ID)

	var test struct {
		Data []models.Invoice `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&test)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Debug(fmt.Sprintf("[GetBills] Received %#v", test.Data))
	return test.Data, nil
}

// Get list of recurring invoices for client
func (inv *InvoiceManager) GetRecurringInvoices(client *models.Client) ([]models.RecurringInvoice, error) {
	res, err := inv.c.Get(inv.endpoint + "/recurring_invoices?status=active&client_id=" + client.ID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var test struct {
		Data []models.RecurringInvoice `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&test)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Debug(fmt.Sprintf("[GetRecurringInvoices] Received %#v", test.Data))
	return test.Data, nil
}
