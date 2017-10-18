package main

import (
	"fmt"
	"net/http"

	doryBackend "github.com/abhishekkr/dory/doryBackend"

	golenv "github.com/abhishekkr/gol/golenv"
	gin "github.com/gin-gonic/gin"
)

var (
	HTTPAt     = golenv.OverrideIfEnv("DORY_HTTP", ":8080")
	VaultAddr  = golenv.OverrideIfEnv("VAULT_ADDR", "http://127.0.0.1:8200")
	VaultToken = golenv.OverrideIfEnv("VAULT_TOKEN", "configure-env-var-VAULT_TOKEN")
)

func main() {
	GinUp(HTTPAt)
	fmt.Println("bye .")
}

func doryHelp(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"help.html",
		gin.H{"title": "Help"},
	)
}

func ginCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}

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

func GinUp(listenAt string) {
	vault := doryBackend.NewVault(VaultAddr, VaultToken)
	localAuth := doryBackend.NewLocalAuth("dory")

	router := gin.Default()
	router.Use(ginHandleErrors)
	router.Use(ginCors())
	router.LoadHTMLGlob("templates/*")
	router.Static("/images", "w3assets/images")

	router.GET("/help", doryHelp)

	v_0_1 := router.Group("/v0.1")
	{
		v_0_1.GET("/vault", vault.AuthList)
		v_0_1.GET("/vault/:uuid", vault.Get)
		v_0_1.POST("/vault/:uuid", vault.AuthMount)
		v_0_1.DELETE("/vault/:uuid", vault.AuthUnmount)

		v_0_1.GET("/local-auth", localAuth.AuthList)
		v_0_1.GET("/local-auth/:uuid", localAuth.Get)
		v_0_1.POST("/local-auth/:uuid", localAuth.AuthMount)
		v_0_1.DELETE("/local-auth/:uuid", localAuth.AuthUnmount)
	}

	router.Run(listenAt)
}
