package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Result interface{} `json:"result,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func GetResource(c *gin.Context) (string, io.ReadCloser, error) {
	var target string
	var body io.ReadCloser
	var err error
	switch c.Request.Method {
	case "GET":
		target = "http://localhost:8080/schedules" + c.Request.URL.Path + "?where=" + c.Query("where")
		body = nil; break
	case "POST":
		target = "http://localhost:8081/schedules" + c.Request.URL.Path
		body = c.Request.Body; break
	case "DELETE":
		target = "http://localhost:8083/schedules" + c.Request.URL.Path
		body = nil; break
	case "PUT":
		target = "http://localhost:8082/schedules" + c.Request.URL.Path
		body = c.Request.Body; break
	default:
		err = fmt.Errorf("Method %s not supported", c.Request.Method)
	}	
	return target, body, err
}

func Proxy(c *gin.Context) {
	tr := &http.Transport{
		DisableKeepAlives: true,
		DisableCompression: true,
	}
	target, body, err := GetResource(c); if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(c.Request.Method, target, body); if err != nil {		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	resp, err := client.Do(req); if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	respBody, _ := io.ReadAll(resp.Body)
	var a response
	json.Unmarshal(respBody, &a)
	c.JSON(http.StatusOK, a)
}

func main() {
	r := gin.New()
	r.Any("/*any", Proxy)
	r.Run(":80")
}