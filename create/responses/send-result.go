package responses

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message interface{} `json:"message"`
}

func (e *SuccessResponse) Send(res http.ResponseWriter, status int) {
	res.Header().Add("content-type", "application/json")
	res.WriteHeader(status)
	if err := json.NewEncoder(res).Encode(e); err != nil {
		res.Write([]byte("Something went wrong with your response. Please verify if your entity was created successfully"))
	}
}