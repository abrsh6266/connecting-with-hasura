package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Allow requests from this origin
	r.Use(cors.New(config))

	r.POST("/try", GetResult)
	r.Run(":8082")
}

func GetResult(ctx *gin.Context) {
	type Material struct {
		Name     	 string `json:"name"`
		Id       	 int `json:"Id"`
		Model 	     string `json:"model"`
		Processor 	 string `json:"processor"`
	}
	
	var result struct {
		Data struct {
			Materials []Material `json:"Material"`
		} `json:"data"`
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	respBody, err := HasuraRequest(string(body))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect Email"})
		return
	}
	
	if json.Unmarshal(respBody,&result); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse category data"})
		return
	}
	ctx.JSON(http.StatusOK,result)
}

var client = &http.Client{}

func HasuraRequest(query string) ([]byte, error) {
	method := "POST"
	req, err := http.NewRequest(method, "https://musical-mink-40.hasura.app/v1/graphql", strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-hasura-admin-secret", "RQxulF8331e3eg6pAq05ikjfY7YsQvnr1wwjmclMbJoagm0IX9k1t9P3m93vfNI6")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}