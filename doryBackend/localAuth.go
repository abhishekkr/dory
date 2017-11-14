package doryBackend

/*
Local Auth Backend for Dory
*/

import (
	"io/ioutil"
	"net/http"
	"strconv"

	doryMemory "github.com/abhishekkr/dory/doryMemory"

	"github.com/gin-gonic/gin"
)

/*
LocalAuth is a struct to maintain connection details for a Local-Auth and single item construct for actions.
*/
type LocalAuth struct {
	Cache doryMemory.DataStore
	Disk  doryMemory.DataStore
	Item  doryMemory.LocalAuth
}

/*
NewLocalAuth instantiates and return a LocalAuth struct in reference to any usable Vault backend.
*/
func NewLocalAuth(cacheName string) LocalAuth {
	localAuth := LocalAuth{
		Cache: doryMemory.NewLocalAuthStore(cacheName),
		Disk:  doryMemory.NewDiskv(cacheName),
		Item:  doryMemory.LocalAuth{},
	}
	return localAuth
}

/*
Get fetchs required auth mapped secret from Local-Auth backend.
*/
func (localAuth LocalAuth) Get(ctx *gin.Context) {
	var datastore doryMemory.DataStore
	if ctx.DefaultQuery("persist", "false") == "false" {
		datastore = localAuth.Cache
	} else {
		datastore = localAuth.Disk
	}

	localAuthItem := localAuth.Item

	localAuthItem.Name = ctx.Param("uuid")
	localAuthItem.Value.Key = []byte(ctx.Request.Header.Get("X-DORY-TOKEN"))

	if !localAuthItem.Get(datastore) {
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.JSON(500, ExitResponse{Msg: "get for required auth identifier failed"})
		return
	}

	response := localAuthItem.Value.DataBlob

	if ctx.DefaultQuery("keep", "false") == "false" {
		if !localAuthItem.Delete(datastore) {
			ctx.JSON(500, ExitResponse{Msg: "auth identifier purge failed", Data: response})
			return
		}
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(response)
}

/*
AuthMount stores a secret mapped with a new auth-path only at Local-Auth with unique auth-token.
*/
func (localAuth LocalAuth) AuthMount(ctx *gin.Context) {
	var datastore doryMemory.DataStore
	if ctx.DefaultQuery("persist", "false") == "false" {
		datastore = localAuth.Cache
	} else {
		datastore = localAuth.Disk
	}

	localAuthItem := localAuth.Item
	localAuthItem.Name = ctx.Param("uuid")

	if localAuthItem.Exists(datastore) {
		ctx.JSON(409, ExitResponse{Msg: "auth identifier conflict"})
		return
	}

	ttlsecond, err := strconv.Atoi(ctx.DefaultQuery("ttlsecond", "0"))
	if err != nil {
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.JSON(400, ExitResponse{Msg: err.Error()})
		return
	}
	localAuthItem.TTLSecond = uint64(ttlsecond)

	localAuthItem.Value.DataBlob, err = ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.JSON(400, ExitResponse{Msg: err.Error()})
		return
	}

	if !localAuthItem.Set(datastore) {
		ctx.JSON(500, ExitResponse{Msg: "auth identifier creation failed"})
		return
	}

	ctx.String(http.StatusOK, string(localAuthItem.Value.Key))
}

/*
AuthUnmount purges a previously local-auth stored mapped to a auth-path if not yet purged by TTL.
*/
func (localAuth LocalAuth) AuthUnmount(ctx *gin.Context) {
	var datastore doryMemory.DataStore
	if ctx.DefaultQuery("persist", "false") == "false" {
		datastore = localAuth.Cache
	} else {
		datastore = localAuth.Disk
	}

	ctx.Writer.Header().Add("Content-Type", "application/json")

	localAuthItem := localAuth.Item
	localAuthItem.Name = ctx.Param("uuid")
	localAuthItem.Value.Key = []byte(ctx.Request.Header.Get("X-DORY-TOKEN"))

	if !localAuthItem.Delete(datastore) {
		ctx.JSON(500, ExitResponse{Msg: "auth identifier purge failed"})
		return
	}

	ctx.JSON(200, ExitResponse{Msg: "success"})
}
