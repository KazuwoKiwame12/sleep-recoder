package main

import (
	"net/http"
	"os"
	"time"

	"wake-time/database"

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
	err := client.Save(now, userID)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func main() {
	lambda.Start(handler)
}
