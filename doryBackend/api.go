package doryBackend

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

/*
A struct to create custom JSON HTTP response.
*/
type ExitResponse struct {
	Msg  string `json:"exit-message"`
	Data []byte `json:"data,omitempty"`
}

/*
To return []byte mapped json response of ExitResponse/
*/
func (response ExitResponse) JSON() (jsonResponse []byte) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		jsonResponse, _ = json.Marshal("{\"error\": \"exit response generation failed\"}")
	}
	return
}

/*
For response handling at API Paths yet WIP.
*/
func wip(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "application/json")

	response := ExitResponse{Msg: "WIP"}
	ctx.JSON(200, response)
}
