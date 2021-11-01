package main

import (
	"fmt"
	"log"
	"math"
	"notify/bucket"
	"notify/db"
	"notify/utility"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/line/line-bot-sdk-go/linebot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// TODO log.Fatalにするか、Printにするかを要考慮...止まらずに次の処理をして欲しいときはPrintに修正する
func main() {
	// line_botの作成
	bot, err := linebot.New(
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
	nowJST := utility.CreateDateWIthJst()
	srs, err := client.ListInWeekForAllUser(nowJST)
	if err != nil {
		log.Fatal(err)
	}
	// 全てのuserIDの抽出
	userIDs := srs.RetrieveUserIDs() //昇順データ
	log.Printf("userIDs: %v\n", userIDs)

	// 上記の睡眠記録データを用いてグラフの作成
	for _, id := range userIDs {
		data := []plotter.Values{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}} // ユーザ一人のグラフデータ
		numOfSrs := 0
		for i, sr := range srs { //srsはuserIDとtimeWの昇順で並んでいる
			if id != sr.UserID {
				numOfSrs = i
				break
			}
			index := utility.GetDiffOfDays(utility.CreateDateWithUnix(sr.TimeW), nowJST.AddDate(0, 0, -6))

			wakeDate := utility.CreateDateWithUnix(sr.TimeW)
			wakeStartDate := utility.CreateStartDate(wakeDate.Year(), wakeDate.Month(), wakeDate.Day())
			bedinHour := utility.CreateDateWithUnix(sr.TimeB).Sub(wakeStartDate).Hours()
			wakeHour := wakeDate.Sub(wakeStartDate).Hours()
			data[index] = plotter.Values{
				math.Round(bedinHour*10) / 10,
				math.Round(wakeHour*10) / 10,
			}
		}
		srs = srs[numOfSrs:] // 取得したデータ数削除する
		log.Printf("data[%s]: %v\n", id, data)
		// グラフの画像を作成する
		if err := createPlotImage(data); err != nil {
			log.Printf("error(createPlotImage): %v\n", err)
			msg := linebot.NewTextMessage("睡眠記録のグラフ化ができませんでした")
			if _, err := bot.PushMessage(id, msg).Do(); err != nil {
				log.Printf("error(pushMessage at 74): %v\n", err)
			}
			continue
		}
		// s3にuploadする
		sessForS3 := session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-northeast-3")}))
		uploader := bucket.NewImageUploader(sessForS3, os.Getenv("BUCKET_NAME"), fmt.Sprintf("sleeprecord-plot_%s_%d-%v-%d.png", id, nowJST.Year(), nowJST.Month(), nowJST.Day()))
		url, err := uploader.UploadImage("sleeprecord-plot.png")
		if err != nil {
			log.Printf("error(UploadImage): %v\n", err)
			msg := linebot.NewTextMessage("システム内でエラーが発生しました")
			if _, err := bot.PushMessage(id, msg).Do(); err != nil {
				log.Printf("error(pushMessage at 86): %v\n", err)
			}
			continue
		}
		// lineに通知
		msg := linebot.NewImageMessage(url, url)
		if _, err := bot.PushMessage(id, msg).Do(); err != nil {
			log.Printf("error(pushMessage at 94): %v\n", err)
			continue
		}
	}
}

func createPlotImage(data []plotter.Values) error {
	p := plot.New()

	p.Title.Text = "Sleep Record for a week"
	p.Y.Label.Text = "time(based AM0:00)"

	now := utility.CreateDateWIthJst()
	if err := plotutil.AddBoxPlots(p, vg.Points(20),
		fmt.Sprintf("%d/%d", int(now.AddDate(0, 0, -6).Month()), now.AddDate(0, 0, -6).Day()), data[6],
		fmt.Sprintf("%d/%d", int(now.AddDate(0, 0, -5).Month()), now.AddDate(0, 0, -5).Day()), data[5],
		fmt.Sprintf("%d/%d", int(now.AddDate(0, 0, -4).Month()), now.AddDate(0, 0, -4).Day()), data[4],
		fmt.Sprintf("%d/%d", int(now.AddDate(0, 0, -3).Month()), now.AddDate(0, 0, -3).Day()), data[3],
		fmt.Sprintf("%d/%d", int(now.AddDate(0, 0, -2).Month()), now.AddDate(0, 0, -2).Day()), data[2],
		fmt.Sprintf("%d/%d", int(now.AddDate(0, 0, -1).Month()), now.AddDate(0, 0, -1).Day()), data[1],
		fmt.Sprintf("%d/%d", int(now.Month()), now.Day()), data[0],
	); err != nil {
		return err
	}

	if err := p.Save(12*vg.Inch, 7*vg.Inch, "sleeprecord-plot.png"); err != nil {
		return err
	}
	return nil
}
