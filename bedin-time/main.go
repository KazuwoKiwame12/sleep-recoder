package main

import (
	// "errors"

	"net/http"
	"os"
	"time"

	"bedin-time/database"
	"bedin-time/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// time validation
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	if err := utility.ValidateBedintime(now.Hour()); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Body:       err.Error(),
		}, nil
	}

	// Request
	userID := request.PathParameters["user_id"]

	// DynamoDB
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-3").WithEndpoint(os.Getenv("DYNAMODB_ENDPOINT"))
	client := database.NewClient(sess, config)

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
