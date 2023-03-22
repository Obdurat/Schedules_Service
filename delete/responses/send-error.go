package responses

import (
	"encoding/json"
	"net/http"
)

type ErrorRespose struct {
	Error interface{} `json:"error"`
}

func (e *ErrorRespose) Send(res http.ResponseWriter, status int) {
	res.Header().Add("content-type", "application/json")
	res.WriteHeader(status)
	if err := json.NewEncoder(res).Encode(e); err != nil {
		res.Write([]byte("Something is wrong with the server. Please try again later"))
	}
}