package doryBackend

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ExitResponse struct {
	Msg string `json:"exit-message"`
}

func (response ExitResponse) JSON() (jsonResponse []byte) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		jsonResponse, _ = json.Marshal("{\"error\": \"exit response generation failed\"}")
	}
	return
}

func wip(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "application/json")

	response := ExitResponse{Msg: "WIP"}
	ctx.JSON(200, response)
}
