package wgez

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"

	"plumpalbert.xyz/plumpwire/models"
)

type WGEasy struct {
	// Host of wg-easy
	host string

	// List of clients
	Clients []models.WG_Client

	Devices map[string][]models.WG_Client
}

// Create new instance of wg
func New(host string) WGEasy {
	wg := WGEasy{
		host:    host,
		Clients: []models.WG_Client{},
		Devices: map[string][]models.WG_Client{},
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
	if err != nil {
		fmt.Println("Could not retreive list of clients")
		fmt.Println(err.Error())
		return err
	}

	for _, c := range wg.Clients {
		re, err := regexp.Compile(`(?P<Username>[\d\w_]+)\s\[(?P<Device>.+)\]`)
		if err != nil {
			continue
		}

		matches := re.FindStringSubmatch(c.Name)
		if matches != nil {
			c.Username = matches[re.SubexpIndex("Username")]
			c.DeviceName = matches[re.SubexpIndex("Device")]

			if wg.Devices[c.Username] == nil {
				wg.Devices[c.Username] = []models.WG_Client{}
			}

			idx := slices.IndexFunc(wg.Devices[c.Username], func(o models.WG_Client) bool {
				return c.ID == o.ID
			})

			if idx == -1 {
				wg.Devices[c.Username] = append(wg.Devices[c.Username], c)
				continue
			}

			wg.Devices[c.Username][idx] = c
		}
	}
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
func (wg WGEasy) GetDevices(user_name string) ([]models.WG_Client, error) {
	err := wg.GetClients()
	if err != nil {
		return nil, err
	}
	return wg.Devices[user_name], nil
}
