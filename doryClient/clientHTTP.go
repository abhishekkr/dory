package doryClient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/abhishekkr/gol/golerror"
	"github.com/abhishekkr/gol/golhttpclient"
	"github.com/abhishekkr/gol/golrandom"
)

type DoryClient struct {
	BaseUrl       string
	Backend       string //local-auth
	Key           string
	Value         []byte
	Token         string
	KeyTTL        int //seconds, usable in cache mode only i.e. when persist is false
	Persist       bool
	ReadNotDelete bool
}

func (dory *DoryClient) userBackend() string {
	backend := dory.Backend
	if backend == "" {
		backend = "local-auth"
	}
	return backend
}

func (dory *DoryClient) adminBackend() string {
	if dory.Persist {
		return "disk"
	}
	return "cache"
}

func (dory *DoryClient) httpUserUrl(request *golhttpclient.HTTPRequest) {
	request.Url = fmt.Sprintf("%s/%s/%s", dory.BaseUrl, dory.userBackend(), dory.Key)
}

func (dory *DoryClient) httpAdminUrl(request *golhttpclient.HTTPRequest) {
	request.Url = fmt.Sprintf("%s/admin/store/%s", dory.BaseUrl, dory.adminBackend())
	if dory.Key != "" {
		request.Url = fmt.Sprintf("%s/%s", request.Url, dory.Key)
	}
}

func (dory *DoryClient) httpUserHeaders(request *golhttpclient.HTTPRequest) {
	request.HTTPHeaders = map[string]string{
		"X-DORY-TOKEN": dory.Token,
	}
}

func (dory *DoryClient) httpAdminHeaders(request *golhttpclient.HTTPRequest) {
	request.HTTPHeaders = map[string]string{
		"X-DORY-ADMIN-TOKEN": dory.Token,
	}
}

func (dory *DoryClient) httpParams(request *golhttpclient.HTTPRequest) {
	ttlsecond := strconv.Itoa(dory.KeyTTL)
	if ttlsecond == "" {
		ttlsecond = "300"
	}
	request.GetParams = map[string]string{
		"keep":      fmt.Sprintf("%t", dory.ReadNotDelete),
		"persist":   fmt.Sprintf("%t", dory.Persist),
		"ttlsecond": ttlsecond,
	}
}

func (dory *DoryClient) httpUserRequest() golhttpclient.HTTPRequest {
	request := golhttpclient.HTTPRequest{}
	dory.httpUserUrl(&request)
	dory.httpUserHeaders(&request)
	dory.httpParams(&request)
	return request
}

func (dory *DoryClient) httpAdminRequest() golhttpclient.HTTPRequest {
	request := golhttpclient.HTTPRequest{}
	dory.httpAdminUrl(&request)
	dory.httpAdminHeaders(&request)
	dory.httpParams(&request)
	return request
}

func (dory *DoryClient) Set() (err error) {
	if dory.Key == "" {
		dory.Key = fmt.Sprintf("dory-%s", golrandom.Token(10))
	}

	request := dory.httpUserRequest()
	request.Body = bytes.NewBuffer(dory.Value)

	dory.Token, err = request.Post()
	return
}

func (dory *DoryClient) ShareSecret(value []byte) (err error) {
	dory.Value = value
	err = dory.Set()
	return
}

func (dory *DoryClient) ShareSecretFromFile(filepath string) (err error) {
	requestBody, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}

	err = dory.ShareSecret(requestBody)
	return
}

func (dory *DoryClient) Get() (err error) {
	if dory.BaseUrl == "" {
		err = golerror.Error(123, "dory url can't be empty")
		return
	}
	if dory.Key == "" || dory.Token == "" {
		err = golerror.Error(123, "key or token can't be empty")
		return
	}

	request := dory.httpUserRequest()

	response, err := request.Get()
	dory.Value = []byte(response)
	return
}

func (dory *DoryClient) FetchSecret() (value []byte, err error) {
	err = dory.Get()
	value = dory.Value
	return
}

func (dory *DoryClient) RefreshSecret() (value []byte, err error) {
	readNotDelete := dory.ReadNotDelete
	dory.ReadNotDelete = true
	value, err = dory.FetchSecret()
	dory.ReadNotDelete = readNotDelete
	return
}

func (dory *DoryClient) Del() (err error) {
	if dory.BaseUrl == "" {
		err = golerror.Error(123, "dory url can't be empty")
		return
	}
	if dory.Key == "" || dory.Token == "" {
		err = golerror.Error(123, "key and token required to purge")
		return
	}

	request := dory.httpUserRequest()

	response, err := request.Get()
	dory.Value = []byte(response)
	return
}

func (dory *DoryClient) PurgeSecret() (err error) {
	return dory.Del()
}

func (dory *DoryClient) PurgeAll() (err error) {
	if dory.BaseUrl == "" {
		err = golerror.Error(123, "dory url can't be empty")
		return
	}
	if dory.Token == "" {
		err = golerror.Error(123, "admin token required to purge")
		return
	}

	request := dory.httpAdminRequest()

	response, err := request.Delete()
	dory.Value = []byte(response)
	return
}

func (dory *DoryClient) PurgeOne() (err error) {
	if dory.BaseUrl == "" {
		err = golerror.Error(123, "dory url can't be empty")
		return
	}
	if dory.Token == "" {
		err = golerror.Error(123, "admin token required to purge")
		return
	}

	request := dory.httpAdminRequest()

	response, err := request.Delete()
	dory.Value = []byte(response)
	return
}

func (dory *DoryClient) List() (err error) {
	if dory.BaseUrl == "" {
		err = golerror.Error(123, "dory url can't be empty")
		return
	}
	if dory.Token == "" {
		err = golerror.Error(123, "admin token required to purge")
		return
	}

	request := dory.httpAdminRequest()

	response, err := request.Get()
	dory.Value = []byte(response)
	return
}

func (dory *DoryClient) Ping() (err error) {
	if dory.BaseUrl == "" {
		err = golerror.Error(123, "dory url can't be empty")
		return
	}
	request := golhttpclient.HTTPRequest{}
	request.Url = fmt.Sprintf("%s/ping", dory.BaseUrl)

	response, err := request.Get()
	dory.Value = []byte(response)
	return
}
