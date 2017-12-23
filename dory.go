package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	doryBackend "github.com/abhishekkr/dory/doryBackend"

	golcrypt "github.com/abhishekkr/gol/golcrypt"
	golenv "github.com/abhishekkr/gol/golenv"
	gollog "github.com/abhishekkr/gol/gollog"

	gin "github.com/gin-gonic/gin"
)

var (
	/*
		HTTPAt specifies server's listen-at config, can be overridden by env var DORY_HTTP. Defaults to '':8080'.
	*/
	HTTPAt = golenv.OverrideIfEnv("DORY_HTTP", ":8080")

	doryMode      = flag.String("mode", "server", "run mode, allowed modes are client and server, defaults server")
	doryUrl       = flag.String("url", "", "url for dory server to be talked to")
	doryKey       = flag.String("key", "", "key name to be provided to dory")
	doryVal       = flag.String("val", "", "value to be provided to dory, required when trying to Post or Delete a key")
	doryValFile   = flag.String("val-from", "", "value from a file to be provided to dory, required when trying to Post or Delete a key")
	doryToken     = flag.String("token", "", "token for secret, required when trying to Get or Delete a key")
	doryClientAxn = flag.String("task", "ping", "the kind of action dory client need to perform, supports {set,get,del,list,purge,ping}; defaults ping")
	doryKeyTTL    = flag.Int("ttl", 300, "ttl for key, if it's set task for cache datastore; defaults 300 sec")
	doryPersist   = flag.Bool("persist", false, "to decide datastore as cache or disk, defaults as false for cache")
	doryReadNKeep = flag.Bool("keep", false, "to decide if to purge key post read or not, defaults as false for purge on read")
)

func main() {
	flag.Parse()

	if *doryMode == "server" {
		gollog.Debug("starting dory as server")
		ginUp(HTTPAt)
	} else if *doryMode == "client" {
		doryClient()
	} else {
		gollog.Err(fmt.Sprintf("wrong run mode '%s' passed to dory", *doryMode))
	}
	gollog.Debug("bye .")
}

/*
ginCors to set required HTTP configs.
*/
func ginCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}

/*
ginHandleErrors to manage issues at server side.
*/
func ginHandleErrors(ctx *gin.Context) {
	ctx.Next()
	errorToPrint := ctx.Errors.ByType(gin.ErrorTypePublic).Last()
	if errorToPrint != nil {
		ctx.JSON(500, gin.H{
			"status":  500,
			"message": errorToPrint.Error(),
		})
	}
}

/*
ginUp maps all routing logic and starts server.
*/
func ginUp(listenAt string) {
	localAuth := doryBackend.NewLocalAuth("dory")

	router := gin.Default()
	router.Use(ginHandleErrors)
	router.Use(gollog.GinLogrus(), gin.Recovery())
	router.Use(ginCors())
	router.LoadHTMLGlob("templates/*")
	router.Static("/images", "w3assets/images")
	router.StaticFile("/favicon.ico", "w3assets/favicon.ico")

	router.GET("/help", doryBackend.DoryHelp)

	router.GET("/ping", localAuth.DoryPing)

	router.GET("/local-auth/:uuid", localAuth.Get)
	router.POST("/local-auth/:uuid", localAuth.AuthMount)
	router.DELETE("/local-auth/:uuid", localAuth.AuthUnmount)

	router.GET("/admin/store/:datastore", localAuth.List)
	router.DELETE("/admin/store/:datastore", localAuth.Purge)

	router.Run(listenAt)
}

/*
doryClient handles calling dory from commandline
*/
func doryClient() {
	var err error

	goldory := golcrypt.Dory{
		BaseUrl:       *doryUrl,
		Key:           *doryKey,
		Token:         *doryToken,
		KeyTTL:        *doryKeyTTL,
		Persist:       *doryPersist,
		ReadNotDelete: *doryReadNKeep,
	}

	goldory.Value = []byte(*doryVal)
	_, err = os.Stat(*doryValFile)
	if *doryValFile != "" && err == nil {
		goldory.Value, err = ioutil.ReadFile(*doryValFile)
		if err != nil {
			gollog.Err(err.Error())
			return
		}
	}

	doryBackend.HandleClientAuth(goldory, *doryClientAxn)
}
