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
type ResponseData struct {
	Status int
	Output []byte
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Body: ", request.Body)
	fmt.Println("PATH: ", request.Path)
	resource := request.Resource
	responseData := &ResponseData{}

	if resource == "/mutant" && request.HTTPMethod == "POST" {
		req := &Request{}
		if err := json.Unmarshal([]byte(request.Body), req); err != nil {
			panic(err)
		}
		data := models.IsMutant(req.Adn)
		rest, _ := json.Marshal(data)
		responseData.Output = rest
		responseData.Status = data.Status
	} else if resource == "/stats" && request.HTTPMethod == "GET" {
		data := models.GetDataCollection()
		rest, _ := json.Marshal(data)
		responseData.Output = rest
		responseData.Status = http.StatusOK
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseData.Output),
		StatusCode: responseData.Status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
