package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"plumpalbert.xyz/plumpwire/models"
)

type WGEasy struct {
	// Host of wg-easy
	host string

	// List of clients
	Clients []models.WG_Client
}

var wg_host = os.Getenv("WG_HOST")

// Get list of wg-easy clients
func (wg WGEasy) GetClients() error {
	res, err := http.Get(wg.host + `/api/wireguard/client/`)

	if err != nil {
		fmt.Println("Could not retreive list of clients")
		fmt.Println(err.Error())
		return err
	}

	err = json.NewDecoder(res.Body).Decode(&wg.Clients)
	return nil
}

// Get client configuration file
func (wg WGEasy) GetClientConfig(client_id string) ([]byte, error) {
	res, err := http.Get(wg_host + `/api/wireguard/client/` + client_id + `/configuration`)

	if err != nil {
		fmt.Println("Could not get configuration")
		fmt.Println(err.Error())
		return nil, err
	}

	return io.ReadAll(res.Body)
}
