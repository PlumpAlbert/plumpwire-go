package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var wg_host = os.Getenv("WG_HOST")

func get_config(id string) ([]byte, error) {
	if wg_host == "" {
		panic("WG_HOST is not defined")
	}
	res, err := http.Get(wg_host + `/api/wireguard/client/` + id + `/configuration`)

	if err != nil {
		fmt.Println("Could not get configuration")
		fmt.Println(err.Error())
		return nil, err
	}

	return io.ReadAll(res.Body)
}
