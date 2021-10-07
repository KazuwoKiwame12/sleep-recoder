package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sleep-manager/db"
	"sleep-manager/usecase"
	"sleep-manager/utility"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	ErrorVerify = errors.New("error: failed to verify signature")
	ErrorInit   = errors.New("error: failed to create bot")
)

type input struct {
	command string
	userID  string
	getter  usecase.Getter
	setter  usecase.Setter
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	secret := os.Getenv("CHANNEL_SECRET")
	if !verifySignature(secret, request.Headers["x-line-signature"], []byte(request.Body)) {
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

	bot, err := linebot.New(
		secret,
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       ErrorInit.Error(),
		}, ErrorInit
	}

	endPoint := os.Getenv("DYNAMODB_ENDPOINT")
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-3").WithEndpoint(endPoint)
	client := db.NewSleepRecordClient(tableName, sess, config)
	getter := usecase.Getter{C: client}
	setter := usecase.Setter{C: client}
	for _, event := range content.Events {
		if event.Type == linebot.EventTypeMessage {
			switch reqMsg := event.Message.(type) {
			case *linebot.TextMessage:
				resMsg := execCommand(
					input{
						command: reqMsg.Text,
						userID:  event.Source.UserID,
						setter:  setter,
						getter:  getter,
					},
				)
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(resMsg)).Do(); err != nil {
					return events.APIGatewayProxyResponse{}, errors.New("failed: can't reply message")
				}
			default:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(utility.MessageDefault)).Do(); err != nil {
					return events.APIGatewayProxyResponse{}, errors.New("failed: can't reply message")
				}
			}
		}
	}

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

func execCommand(i input) string {
	switch utility.ValidateCommand(i.command) {
	case utility.CommandBedin:
		if err := i.setter.SaveBedinTime(i.userID); err != nil {
			return utility.MessageError
		}
		return utility.MessageSuccessRecord
	case utility.CommandWake:
		if err := i.setter.SaveWakeTime(i.userID); err != nil {
			return utility.MessageError
		}
		return utility.MessageSuccessRecord
	case utility.CommandFiveDays:
		msg, err := i.getter.ListRecordsInFiveDays(i.userID)
		if err != nil {
			return utility.MessageError
		}
		if len(msg.Record) == 0 {
			return utility.MessageNotFound
		}
		msgJson, err := json.Marshal(&msg)
		if err != nil {
			return utility.MessageError
		}
		return string(msgJson)
	case utility.CommandMonth:
		slice := strings.Split(i.command, " ")
		year, err := strconv.Atoi(slice[1])
		if err != nil {
			return utility.MessageError
		}
		month, err := strconv.Atoi(slice[2])
		if err != nil {
			return utility.MessageError
		}
		msg, err := i.getter.ListRecordsInMonth(year, time.Month(month), i.userID)
		if err != nil {
			return utility.MessageError
		}
		if len(msg.Record) == 0 {
			return utility.MessageNotFound
		}
		msgJson, err := json.Marshal(&msg)
		if err != nil {
			return utility.MessageError
		}
		return string(msgJson)
	case utility.CommandHelp:
		return utility.MessageHelp
	case utility.CommandDefault:
		return utility.MessageDefault
	}
	return utility.MessageError
}

func main() {
	lambda.Start(handler)
}
