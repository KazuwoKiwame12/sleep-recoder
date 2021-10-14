package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
	"golang.org/x/text/width"
)

type input struct {
	command string
	userID  string
	getter  usecase.Getter
	setter  usecase.Setter
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	secret := os.Getenv("LINE_CHANNEL_SECRET")
	if !verifySignature(secret, request.Headers["x-line-signature"], []byte(request.Body)) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       utility.ErrorVerify.Error(),
		}, utility.ErrorVerify
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
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       utility.ErrorInit.Error(),
		}, utility.ErrorInit
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-3")
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
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       utility.ErrorReply.Error(),
					}, utility.ErrorReply
				}
			default:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(utility.MessageDefault)).Do(); err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       utility.ErrorReply.Error(),
					}, utility.ErrorReply
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
			if err == utility.ErrorBedinTime {
				return utility.MessageBedinTimeError
			}
			return utility.MessageSystemError
		}
		return utility.MessageSuccessRecord
	case utility.CommandWake:
		if err := i.setter.SaveWakeTime(i.userID); err != nil {
			return utility.MessageNotSleep
		}
		return utility.MessageSuccessRecord
	case utility.CommandFiveDays:
		msg, err := i.getter.ListRecordsInFiveDays(i.userID)
		if err != nil {
			return utility.MessageSystemError
		}
		if len(msg.Record) == 0 {
			return utility.MessageNotFound
		}
		return msg.GetLineMessage()
	case utility.CommandMonth:
		slice := strings.Split(width.Narrow.String(i.command), " ")
		year, err := strconv.Atoi(slice[1])
		if err != nil {
			return utility.MessageNotExistYear
		}
		month, err := strconv.Atoi(slice[2])
		if err != nil {
			return utility.MessageNotExistMonth
		}
		m := time.Month(month)
		if m < time.January || m > time.December {
			return utility.MessageNotExistMonth
		}
		msg, err := i.getter.ListRecordsInMonth(year, m, i.userID)
		if err != nil {
			return utility.MessageSystemError
		}
		if len(msg.Record) == 0 {
			return utility.MessageNotFound
		}
		return msg.GetLineMessage()
	case utility.CommandHelp:
		return utility.MessageHelp
	case utility.CommandDefault:
		return utility.MessageDefault
	}
	return utility.MessageSystemError
}

func main() {
	lambda.Start(handler)
}
