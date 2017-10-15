package doryBackend

/*
Vault Backend for Dory
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/abhishekkr/gol/golhttpclient"
)

type Vault struct {
	BaseAddr  string
	AuthToken string
	Request   golhttpclient.HTTPRequest
}

type VaultAuthBackend struct {
	MountPoint  string `json:"mountpoint"`
	Type        string `json:"type"`
	Local       bool   `json:"local"`
	Description string `json:"description"`
}

func NewVault(baseAddr string, authToken string) Vault {
	vaultHTTPHeaders := map[string]string{
		"X-Vault-Token": authToken,
	}

	request := golhttpclient.HTTPRequest{
		HTTPHeaders: vaultHTTPHeaders,
		Url:         baseAddr,
	}

	return Vault{
		BaseAddr:  baseAddr,
		AuthToken: authToken,
		Request:   request,
	}
}

func (vault Vault) AuthList() {
	vault.Request.Url = fmt.Sprintf("%s/v1/sys/auth", vault.Request.Url)
	response, err := vault.Request.Get()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(response)
}

func (vault Vault) AuthMount(auth VaultAuthBackend) {
	vault.Request.Url = fmt.Sprintf("%s/v1/sys/auth/%s", vault.Request.Url, auth.MountPoint)
	vault.Request.HTTPHeaders["Content-Type"] = "application/json"

	requestBodyJson, err := json.Marshal(auth)
	vault.Request.Body = bytes.NewBuffer([]byte(string(requestBodyJson)))
	fmt.Println(vault.Request)

	response, err := vault.Request.Post()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(response)
}

func (vault Vault) AuthUnmount(auth VaultAuthBackend) {
	vault.Request.Url = fmt.Sprintf("%s/v1/sys/auth/%s", vault.Request.Url, auth.MountPoint)

	response, err := vault.Request.Delete()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(response)
}
