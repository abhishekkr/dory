package main

import (
	"fmt"
	"net/http"

	doryBackend "github.com/abhishekkr/dory/doryBackend"

	golenv "github.com/abhishekkr/gol/golenv"
	gin "github.com/gin-gonic/gin"
)

var (
	/*
		HTTPAt specifies server's listen-at config, can be overridden by env var DORY_HTTP. Defaults to '':8080'.
	*/
	HTTPAt = golenv.OverrideIfEnv("DORY_HTTP", ":8080")

	/*
		VaultAddr is Vault Backend URI. Can be overridden by env var VAULT_ADDR. Defaults to 'http://127.0.0.1:8200'.
	*/
	VaultAddr = golenv.OverrideIfEnv("VAULT_ADDR", "http://127.0.0.1:8200")

	/*
		VaultToken is Vault Root Token. Need to be overridden by env var VAULT_TOKEN. Defaults to unusable placeholder value.
	*/
	VaultToken = golenv.OverrideIfEnv("VAULT_TOKEN", "configure-env-var-VAULT_TOKEN")
)

func main() {
	GinUp(HTTPAt)
	fmt.Println("bye .")
}

/*
doryHelp to serve help file for Dory.
*/
func doryHelp(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"help.html",
		gin.H{"title": "Help"},
	)
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
	vault := doryBackend.NewVault(VaultAddr, VaultToken)
	localAuth := doryBackend.NewLocalAuth("dory")

	router := gin.Default()
	router.Use(ginHandleErrors)
	router.Use(ginCors())
	router.LoadHTMLGlob("templates/*")
	router.Static("/images", "w3assets/images")
	router.StaticFile("/favicon.ico", "w3assets/favicon.ico")

	router.GET("/help", doryHelp)

	router.GET("/local-auth/:uuid", localAuth.Get)
	router.POST("/local-auth/:uuid", localAuth.AuthMount)
	router.DELETE("/local-auth/:uuid", localAuth.AuthUnmount)

	alpha := router.Group("/alpha")
	{
		alpha.GET("/vault/:uuid", vault.Get)
		alpha.POST("/vault/:uuid", vault.AuthMount)
		alpha.DELETE("/vault/:uuid", vault.AuthUnmount)
	}

	router.Run(listenAt)
}
