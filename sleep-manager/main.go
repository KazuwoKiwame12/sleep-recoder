package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	ErrorVerify = errors.New("error: failed to verify signature")
	ErrorInit   = errors.New("error: failed to create bot")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	secret := os.Getenv("CHANNEL_SECRET")
	if !verifySignature(secret, request.Headers["x-line-signature"], []byte(request.Body)) {
		log.Println()
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       ErrorVerify.Error(),
		}, ErrorVerify
	}

	content := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal([]byte(request.Body), content); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}

	_, err := linebot.New(
		secret,
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       ErrorInit.Error(),
		}, ErrorInit
	}

	// TODO implement assignig to approproate functions with Line Events

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func verifySignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}

func main() {
	lambda.Start(handler)
}
