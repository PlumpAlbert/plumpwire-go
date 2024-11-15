package wgez

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"plumpalbert.xyz/plumpwire/models"
)

type WGEasy struct {
	// Host of wg-easy
	host string

	// List of clients
	Clients []models.WG_Client
}

// Create new instance of wg
func New(host string) WGEasy {
	wg := WGEasy{
		host:    host,
		Clients: []models.WG_Client{},
	}

	return wg
}

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
	res, err := http.Get(wg.host + `/api/wireguard/client/` + client_id + `/configuration`)

	if err != nil {
		fmt.Println("Could not get configuration")
		fmt.Println(err.Error())
		return nil, err
	}

	return io.ReadAll(res.Body)
}

// Get list of devices per client
func (wg WGEasy) GetDevices(client_name string) ([]models.WG_Client, error) {
	err := wg.GetClients()
	if err != nil {
		return nil, err
	}

	var clients []models.WG_Client

	for _, c := range wg.Clients {
		regexpString := client_name + `\s\[(?P<Device>.+)\]`
		re := regexp.MustCompile(regexpString)

		if re == nil {
			return nil, errors.New(
				"Could not create regular expression from string: \"" +
					client_name + `\s\[(?P<Device>.+)\]"`,
			)
		}

		matches := re.FindStringSubmatch(c.Name)
		if matches != nil {
			fmt.Printf("Matches: %s\n", matches)
			clients = append(clients, c)
		}
	}

	return clients, nil
}
