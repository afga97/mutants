package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afga97/mutants/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Adn []string `json:"adn"`
}

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Output  []byte
}

type ResponseList struct {
	CountMutant int64   `json:"count_mutant_dna"`
	CountHuman  int64   `json:"count_human_dna"`
	Ratio       float32 `json:"ratio"`
	Status      int     `json:"status"`
	Output      []byte
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Body: ", request.Body)
	fmt.Println("PATH: ", request.Path)
	resource := request.Resource
	respon := &Response{}
	responList := &ResponseList{}

	if resource == "/mutant" && request.HTTPMethod == "POST" {
		req := &Request{}
		if err := json.Unmarshal([]byte(request.Body), req); err != nil {
			panic(err)
		}
		isMutant := models.IsMutant(req.Adn)
		if isMutant {
			respon.Message = "You are a mutant"
			respon.Status = http.StatusOK
		} else {
			respon.Message = "You are a human"
			respon.Status = http.StatusForbidden
		}
		rest, _ := json.Marshal(respon)
		respon.Output = rest
	} else if resource == "/stats" && request.HTTPMethod == "GET" {
		data := models.GetDataCollection()
		rest, _ := json.Marshal(data)
		responList.Output = rest
		responList.Status = http.StatusOK
	}

	var dataResponse string
	var status int
	if respon.Output != nil {
		dataResponse = string(respon.Output)
		status = respon.Status
	} else {
		dataResponse = string(responList.Output)
		status = responList.Status
	}

	return events.APIGatewayProxyResponse{
		Body:       dataResponse,
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
