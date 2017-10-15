package doryBackend

import "encoding/json"

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
