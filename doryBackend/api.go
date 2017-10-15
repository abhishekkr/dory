package doryBackend

import "github.com/gin-gonic/gin"

func (vault Vault) API(ctx *gin.Context) {

	newAuth := VaultAuthBackend{
		MountPoint:  "nemo",
		Type:        "userspace",
		Local:       true,
		Description: "login to find nemo",
	}
	vault.AuthList()
	vault.AuthMount(newAuth)
	vault.AuthList()
	vault.AuthUnmount(newAuth)
	vault.AuthList()

}
