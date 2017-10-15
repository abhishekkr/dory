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

func GinUp(listenAt string) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	vault := doryBackend.NewVault(VaultAddr, VaultToken)

	router.Static("/images", "w3assets/images")

	router.GET("/help", doryHelp)

	v_0_1 := router.Group("/v0.1")
	{
		v_0_1.GET("/vault", vault.AuthList)
		v_0_1.POST("/vault/:uuid", vault.AuthMount)
		v_0_1.DELETE("/vault/:uuid", vault.AuthUnmount)
	}

	router.Run(listenAt)
}
