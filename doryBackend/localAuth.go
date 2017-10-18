package doryBackend

/*
Local Auth Backend for Dory
*/

import (
	"io/ioutil"
	"net/http"
	"strconv"

	doryMemory "github.com/abhishekkr/dory/doryMemory"
	"github.com/muesli/cache2go"

	"github.com/gin-gonic/gin"
)

type LocalAuth struct {
	Store *cache2go.CacheTable
	Item  doryMemory.LocalAuth
}

func NewLocalAuth(cacheName string) LocalAuth {
	localAuth := LocalAuth{
		Store: doryMemory.NewLocalAuthStore(cacheName),
		Item:  doryMemory.LocalAuth{},
	}
	return localAuth
}

func (localAuth LocalAuth) AuthList(ctx *gin.Context) {
	wip(ctx)
}

func (localAuth LocalAuth) Get(ctx *gin.Context) {
	localAuthItem := localAuth.Item

	localAuthItem.Name = ctx.Param("uuid")
	localAuthItem.Value.Key = []byte(ctx.Request.Header.Get("X-DORY-TOKEN"))

	if !localAuthItem.Get(localAuth.Store) {
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.JSON(500, ExitResponse{Msg: "get for required auth identifier failed"})
		return
	}

	response := localAuthItem.Value.DataBlob

	if ctx.DefaultQuery("keep", "false") == "false" {
		if !localAuthItem.Delete(localAuth.Store) {
			ctx.JSON(500, ExitResponse{Msg: "auth identifier purge failed", Data: response})
			return
		}
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(response)
}

func (localAuth LocalAuth) AuthMount(ctx *gin.Context) {
	localAuthItem := localAuth.Item
	localAuthItem.Name = ctx.Param("uuid")

	if localAuthItem.Exists(localAuth.Store) {
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

	if !localAuthItem.Set(localAuth.Store) {
		ctx.JSON(500, ExitResponse{Msg: "auth identifier creation failed"})
		return
	}

	ctx.String(http.StatusOK, string(localAuthItem.Value.Key))
}

func (localAuth LocalAuth) AuthUnmount(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "application/json")

	localAuthItem := localAuth.Item
	localAuthItem.Name = ctx.Param("uuid")
	localAuthItem.Value.Key = []byte(ctx.Request.Header.Get("X-DORY-TOKEN"))

	if !localAuthItem.Delete(localAuth.Store) {
		ctx.JSON(500, ExitResponse{Msg: "auth identifier purge failed"})
		return
	}

	ctx.JSON(200, ExitResponse{Msg: "success"})
}
