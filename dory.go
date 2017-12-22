package main

import (
	"fmt"

	doryBackend "github.com/abhishekkr/dory/doryBackend"

	golenv "github.com/abhishekkr/gol/golenv"
	gollog "github.com/abhishekkr/gol/gollog"

	gin "github.com/gin-gonic/gin"
)

var (
	/*
		HTTPAt specifies server's listen-at config, can be overridden by env var DORY_HTTP. Defaults to '':8080'.
	*/
	HTTPAt = golenv.OverrideIfEnv("DORY_HTTP", ":8080")
)

func main() {
	GinUp(HTTPAt)
	fmt.Println("bye .")
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
GinUp maps all routing logic and starts server.
*/
func GinUp(listenAt string) {
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
