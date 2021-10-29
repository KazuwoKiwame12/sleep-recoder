package main

import (
	"log"
	"notify/db"
	"notify/entity"
	"notify/utility"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	// line_botの作成
	lineBot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	// dynamodbにアクセスするインスタンス作成
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-3")
	client := db.NewSleepRecordClient(tableName, sess, config)

	// db処理で1週間分のデータを読み込む
	srs, err := client.ListInWeekForAllUser(utility.CreateDateWIthJst())
	if err != nil {
		log.Fatal(err)
	}
	// 全てのuserIDの抽出
	userIDs := srs.RetrieveUserIDs() //昇順データ

	// 上記の睡眠記録データを用いてグラフの作成
	for _, id := range userIDs {
		data := make([]entity.PlotData, 7) // ユーザ一人のグラフデータ
		numOfSrs := 0
		for i, sr := range srs { //srsはuserIDとtimeWの昇順で並んでいる
			if id != sr.UserID {
				numOfSrs = i
				break
			}
			index := utility.GetDiffOfDays(utility.CreateDateWithUnix(sr.TimeW), utility.CreateDateWIthJst().AddDate(0, 0, -6))

			data[index] = entity.PlotData{
				TimeB: sr.TimeB,
				TimeW: sr.TimeW,
			}
		}
		srs = srs[numOfSrs:] // 取得したデータ数削除する
		// グラフの画像を作成する
		// lineに通知
	}
}
