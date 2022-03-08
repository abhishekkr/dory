package doryBackend

/*
Local Auth Backend for Dory
*/

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	doryMemory "github.com/abhishekkr/dory/doryMemory"

	"github.com/abhishekkr/gol/golenv"
	"github.com/abhishekkr/gol/golerror"
	"github.com/abhishekkr/gol/gollog"

	"github.com/gin-gonic/gin"
)

var (
	DORY_ADMIN_TOKEN = golenv.OverrideIfEnv("DORY_ADMIN_TOKEN", "")
	DORY_TEMP_DIR    = golenv.OverrideIfEnv("DORY_TEMP_DIR", "temp")
)

func init() {
	if _, err := os.Stat(DORY_TEMP_DIR); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(DORY_TEMP_DIR, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

/*
LocalAuth is a struct to maintain connection details for a Local-Auth and single item construct for actions.
*/
type LocalAuth struct {
	Cache doryMemory.DataStore
	Disk  doryMemory.DataStore
	Item  doryMemory.LocalAuth
}

/*
NewLocalAuth instantiates and return a LocalAuth struct in reference to any usable secret backend.
*/
func NewLocalAuth(cacheName string) LocalAuth {
	localAuth := LocalAuth{
		Cache: doryMemory.NewLocalAuthStore(cacheName),
		Disk:  doryMemory.NewDiskv(cacheName),
		Item:  doryMemory.LocalAuth{},
	}
	return localAuth
}

func (localAuth LocalAuth) ctxPersist(ctx *gin.Context) (datastore doryMemory.DataStore) {
	requestURI := ctx.Request.RequestURI
	requestAt := strings.Split(requestURI, "/")[1]
	if requestAt == "local-cache" {
		gollog.Debug(fmt.Sprintf("key '%s' is provided for memory store with expiry", localAuth.Item.Name))
		datastore = localAuth.Cache
		return
	} else if requestAt == "local-disk" {
		gollog.Debug(fmt.Sprintf("key '%s' is provided for long-term disk store", localAuth.Item.Name))
		datastore = localAuth.Disk
		return
	}

	if ctx.DefaultQuery("persist", "false") == "false" {
		gollog.Debug(fmt.Sprintf("key '%s' is provided for memory store with expiry", localAuth.Item.Name))
		datastore = localAuth.Cache
	} else {
		gollog.Debug(fmt.Sprintf("key '%s' is provided for long-term disk store", localAuth.Item.Name))
		datastore = localAuth.Disk
	}
	return
}

func (localAuth LocalAuth) ctxDatastore(ctx *gin.Context) (datastore doryMemory.DataStore, err error) {
	datastoreType := ctx.Param("datastore")
	if datastoreType == "cache" {
		datastore = localAuth.Cache
	} else if datastoreType == "disk" {
		datastore = localAuth.Disk
	} else {
		err = golerror.Error(123, fmt.Sprintf("store %s is not allowed, only 'cache' and 'disk' are allowed"))
	}
	return
}

func (localAuth LocalAuth) ctxAdminToken(ctx *gin.Context) (err error) {
	adminToken := ctx.Request.Header.Get("X-DORY-ADMIN-TOKEN")

	if len(DORY_ADMIN_TOKEN) < 256 {
		err = golerror.Error(123, "configured admin token length is less than 64 chars, not allowed")
		return
	}
	if DORY_ADMIN_TOKEN != adminToken {
		err = golerror.Error(123, "provided admin token doesn't match configured token")
		return
	}
	return
}

/*
Get fetchs required auth mapped secret from Local-Auth backend.
*/
func (localAuth LocalAuth) Get(ctx *gin.Context) {
	localAuth.Item.Name = ctx.Param("uuid")
	localAuth.Item.Key = []byte(ctx.Request.Header.Get("X-DORY-TOKEN"))
	datastore := localAuth.ctxPersist(ctx)

	if localAuth.Item.Name == "" {
		ctx.JSON(500, ExitResponse{Msg: "passed uuid is empty"})
		return
	}
	if !localAuth.Item.Get(datastore) {
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.JSON(500, ExitResponse{Msg: "get for required auth identifier failed"})
		return
	}

	response := localAuth.Item.Value.DataBlob

	if ctx.DefaultQuery("keep", "false") == "false" {
		if !localAuth.Item.Delete(datastore) {
			ctx.JSON(500, ExitResponse{Msg: "auth identifier purge failed", Data: response})
			return
		}
	} else {
		gollog.Debug(fmt.Sprintf("GET - key '%s' is queried to be not purged", localAuth.Item.Name))
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(response)
}

/*
AuthMount stores a secret mapped with a new auth-path only at Local-Auth with unique auth-token.
*/
func (localAuth LocalAuth) AuthMount(ctx *gin.Context) {
	localAuth.Item.Name = ctx.Param("uuid")
	datastore := localAuth.ctxPersist(ctx)

	if localAuth.Item.Name == "" {
		ctx.JSON(500, ExitResponse{Msg: "passed uuid is empty"})
		return
	}
	if localAuth.Item.Exists(datastore) {
		ctx.JSON(409, ExitResponse{Msg: "auth identifier conflict"})
		return
	}

	ttlsecond, err := strconv.Atoi(ctx.DefaultQuery("ttlsecond", "0"))
	if err != nil {
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.JSON(400, ExitResponse{Msg: err.Error()})
		return
	}
	localAuth.Item.TTLSecond = uint64(ttlsecond)

	fileField := ctx.DefaultQuery("file-field", "")
	if fileField == "" {
		localAuth.Item.Value.DataBlob, err = ioutil.ReadAll(ctx.Request.Body)
	} else {
		dst := path.Join(DORY_TEMP_DIR, localAuth.Item.Name)
		bigfile, _ := ctx.FormFile(fileField)
		ctx.SaveUploadedFile(bigfile, dst)
		localAuth.Item.Value.DataBlob, err = ioutil.ReadFile(dst)
		warn := os.Remove(dst)
		if warn != nil {
			gollog.Warnf(fmt.Sprintf("SET - unable to delete %s from filesystem", dst))
		}
	}

	if err != nil {
		gollog.Err(fmt.Sprintf("SET - key '%s' had failure to read it's data", localAuth.Item.Name))
		ctx.Writer.Header().Add("Content-Type", "application/json")
		ctx.JSON(400, ExitResponse{Msg: err.Error()})
		return
	}
	if len(localAuth.Item.Value.DataBlob) == 0 {
		gollog.Err(fmt.Sprintf("SET - key '%s' is provided with empty data", localAuth.Item.Name))
		ctx.JSON(400, ExitResponse{Msg: "empty data blob recieved"})
		return
	}

	if !localAuth.Item.Set(datastore) {
		ctx.JSON(500, ExitResponse{Msg: "auth identifier creation failed"})
		return
	}

	ctx.String(http.StatusOK, string(localAuth.Item.Key))
}

/*
AuthUnmount purges a previously local-auth stored mapped to a auth-path if not yet purged by TTL.
*/
func (localAuth LocalAuth) AuthUnmount(ctx *gin.Context) {
	localAuth.Item.Name = ctx.Param("uuid")
	localAuth.Item.Key = []byte(ctx.Request.Header.Get("X-DORY-TOKEN"))
	datastore := localAuth.ctxPersist(ctx)

	ctx.Writer.Header().Add("Content-Type", "application/json")

	if localAuth.Item.Name == "" {
		ctx.JSON(500, ExitResponse{Msg: "passed uuid is empty"})
		return
	}
	if !localAuth.Item.Delete(datastore) {
		ctx.JSON(500, ExitResponse{Msg: "auth identifier purge failed"})
		return
	}

	ctx.JSON(200, ExitResponse{Msg: "success"})
}

/*
List shows all keys registered with Dory for datatsore enquired.
*/
func (localAuth LocalAuth) List(ctx *gin.Context) {
	var err error
	ctx.Writer.Header().Add("Content-Type", "application/json")

	datastore, err := localAuth.ctxDatastore(ctx)
	if err != nil {
		ctx.JSON(500, ExitResponse{Msg: err.Error()})
		return
	}

	err = localAuth.ctxAdminToken(ctx)
	if err != nil {
		ctx.JSON(500, ExitResponse{Msg: err.Error()})
		return
	}

	ctx.JSON(200, datastore.List())
}

/*
Purge removes all keys from datastore enquired, without decryption required.
*/
func (localAuth LocalAuth) Purge(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "application/json")

	datastore, err := localAuth.ctxDatastore(ctx)
	if err != nil {
		ctx.JSON(500, ExitResponse{Msg: err.Error()})
		return
	}

	err = localAuth.ctxAdminToken(ctx)
	if err != nil {
		ctx.JSON(500, ExitResponse{Msg: err.Error()})
		return
	}

	ctx.JSON(200, datastore.Purge())
}

/*
PurgeOne removes only one provided key from datastore enquired, without decryption required.
*/
func (localAuth LocalAuth) PurgeOne(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "application/json")

	datastore, err := localAuth.ctxDatastore(ctx)
	if err != nil {
		ctx.JSON(500, ExitResponse{Msg: err.Error()})
		return
	}

	err = localAuth.ctxAdminToken(ctx)
	if err != nil {
		ctx.JSON(500, ExitResponse{Msg: err.Error()})
		return
	}

	localAuth.Item.Name = ctx.Param("uuid")

	if localAuth.Item.Name == "" {
		ctx.JSON(500, ExitResponse{Msg: "passed uuid is empty"})
		return
	}
	if datastore.PurgeOne(localAuth.Item.Name) != nil {
		ctx.JSON(500, ExitResponse{Msg: "purge-one failed"})
		return
	}

	ctx.JSON(200, ExitResponse{Msg: "success"})
}

/*
doryPing to return status for Dory
*/
func (localAuth LocalAuth) DoryPing(ctx *gin.Context) {
	ping := map[string]string{
		"keys-in-cache": fmt.Sprintf("%d", localAuth.Cache.Count()),
		"keys-in-disk":  fmt.Sprintf("%d", localAuth.Disk.Count()),
	}

	ctx.JSON(200, ping)
}
