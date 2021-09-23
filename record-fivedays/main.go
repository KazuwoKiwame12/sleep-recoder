package main

import (
	"encoding/json"
	"net/http"
	"os"
	"record-fivedays/database"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Request
	userID := request.PathParameters["user_id"]

	// DynamoDB
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-3").WithEndpoint(os.Getenv("DYNAMODB_ENDPOINT"))
	client := database.NewClient(sess, config)
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	past := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst).AddDate(0, 0, -4)
	rcs, err := client.Get(past, userID)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	sr_bytes, err := json.Marshal(rcs)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(sr_bytes),
	}, nil
}

func main() {
	lambda.Start(handler)
}
