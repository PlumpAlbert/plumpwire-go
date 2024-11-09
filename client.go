package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Server  Server
	Clients map[string]Client
}

type Server struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	IPAddress  string `json:"address"`
}

type Client struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	IPAddress    string `json:"address"`
	PublicKey    string `json:"publicKey"`
	PrivateKey   string `json:"privateKey"`
	PreSharedKey string `json:"preSharedKey"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	Enabled      bool   `json:"enabled"`
}

func parse_json() {
	jsonFile, err := os.Open("wg0.json")
	if err != nil {
		fmt.Println("Could not open wg0.json")
		panic(err)
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(bytes, &config)

	for _, client := range config.Clients {
		fmt.Printf("User: %s\n", client.Name)
		fmt.Printf("IP: %s\n", client.IPAddress)
		fmt.Printf("Enabled: %t\n\n", client.Enabled)
	}
}
