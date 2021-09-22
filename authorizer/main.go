package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	fmt.Printf("token: %s\n", token)
	tokenSplit := strings.Split(token, "")
	fmt.Printf("tokenSplit: %+v\n", tokenSplit)
	var bearerToken string
	if len(tokenSplit) == 2 {
		bearerToken = tokenSplit[1]
	}
	fmt.Printf("bearer: %s\n", bearerToken)
	if bearerToken != os.Getenv("CHANNEL_ACCESS_TOKEN") {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	return generatePolicy("user", "Allow", request.MethodArn), nil
}

func main() {
	lambda.Start(handler)
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	return authResponse
}
