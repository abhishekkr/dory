package doryBackend

/*
Vault Backend for Dory
*/

import (
	"github.com/abhishekkr/gol/golhttpclient"
	"github.com/gin-gonic/gin"
)

/*
Vault is a struct to maintain connection details for a Hashicorp Vault Backend
*/
type Vault struct {
	BaseAddr  string
	AuthToken string
	Request   golhttpclient.HTTPRequest
}

/*
VaultAuthBackend is a struct to hold each secret store that maps to a auth-path.
*/
type VaultAuthBackend struct {
	MountPoint  string `json:"mountpoint"`
	Type        string `json:"type"`
	Local       bool   `json:"local"`
	Description string `json:"description"`
}

/*
NewVault instantiates and return a Vault struct in reference to any usable Vault backend.
*/
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

/*
AuthList lists not-sensitive details on secrets stored at Vault.
*/
func (vault Vault) AuthList(ctx *gin.Context) {
	wip(ctx)
	return
	/*
		vault.Request.Url = fmt.Sprintf("%s/v1/sys/auth", vault.Request.Url)
		response, err := vault.Request.Get()
		if err != nil {
			ctx.JSON(500, ExitResponse{Msg: err.Error()}.JSON())
			return
		}
		ctx.JSON(200, response)
	*/
}

/*
Get fetchs a required auth mapped secret from Vault backend.
*/
func (vault Vault) Get(ctx *gin.Context) {
	wip(ctx)
	return
	/*
		response := ExitResponse{Msg: "WIP"}.JSON()
		ctx.JSON(200, response)
	*/
}

/*
AuthMount stores a secret mapped with a new auth-path only at Vault with unique auth-token.
*/
func (vault Vault) AuthMount(ctx *gin.Context) {
	wip(ctx)
	return
	/*
		mountPoint := ctx.Param("uuid")

		auth := VaultAuthBackend{
			MountPoint:  mountPoint,
			Type:        "userspace",
			Local:       true,
			Description: fmt.Sprintf("login to find %s", mountPoint),
		}

		vault.Request.Url = fmt.Sprintf("%s/v1/sys/auth/%s", vault.Request.Url, auth.MountPoint)
		vault.Request.HTTPHeaders["Content-Type"] = "application/json"

		requestBodyJson, err := json.Marshal(auth)
		vault.Request.Body = bytes.NewBuffer([]byte(string(requestBodyJson)))
		fmt.Println(vault.Request)

		response, err := vault.Request.Post()
		if err != nil {
			ctx.JSON(500, ExitResponse{Msg: err.Error()}.JSON())
			return
		} else if response != "" {
			ctx.JSON(400, response)
			return
		}
		ctx.JSON(200, ExitResponse{Msg: auth.MountPoint}.JSON())
	*/
}

/*
AuthUnmount purges a previously vault stored auth mapped to a auth-path if not yet purged by TTL.
*/
func (vault Vault) AuthUnmount(ctx *gin.Context) {
	wip(ctx)
	return
	/*
		mountPoint := ctx.Param("uuid")

		vault.Request.Url = fmt.Sprintf("%s/v1/sys/auth/%s", vault.Request.Url, mountPoint)

		response, err := vault.Request.Delete()
		if err != nil {
			ctx.JSON(500, ExitResponse{Msg: err.Error()}.JSON())
			return
		}
		ctx.JSON(200, response)
	*/
}
