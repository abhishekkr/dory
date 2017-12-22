package doryBackend

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
ExitResponse is a struct to create custom JSON HTTP response.
*/
type ExitResponse struct {
	Msg  string `json:"exit-message"`
	Data []byte `json:"data,omitempty"`
}

/*
JSON returns []byte mapped json response of ExitResponse/
*/
func (response ExitResponse) JSON() (jsonResponse []byte) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		jsonResponse, _ = json.Marshal("{\"error\": \"exit response generation failed\"}")
	}
	return
}

/*
wip sets response handling at API Paths yet WIP.
*/
func Wip(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "application/json")

	response := ExitResponse{Msg: "WIP"}
	ctx.JSON(200, response)
}

/*
doryHelp to serve help file for Dory.
*/
func DoryHelp(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"help.html",
		gin.H{"title": "Help"},
	)
}
