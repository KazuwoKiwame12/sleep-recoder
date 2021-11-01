# sleep-recorder
このアプリでは，LINEを通して簡単に睡眠時間をお手軽に記録し、1週間・1ヶ月の**睡眠情報**やその**評価**を**簡単に可視化**する

## Badges
[![CI](https://github.com/KazuwoKiwame12/sleep-recorder/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/KazuwoKiwame12/sleep-recorder/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/KazuwoKiwame12/sleep-recorder)](/LICENSE)

## 主機能
1. ユーザがLINE画面で"眠る"を打つと、システムは睡眠時刻を保存する
2. ユーザがLINE画面で"起きた"を打つと、システムは起床時間と1との間の睡眠時間を記録する
3. ユーザがLINE画面で"取得"を打つと、システムは現時点から4日前までの5日間の睡眠記録リストを表示する
4. ユーザがLINE画面で"取得 年(数字) 月(数字)"を打つ、システムは該当する月の睡眠記録リストを表示する
6. システムは1週間に1度、月曜日の夜23時に、LINEにその週の睡眠記録リストのグラフを送信する

## 要素技術
- フロントエンド
    - LINE
- バックエンド
    - サーバ: lambda
    - ゲートウェイ: API Gateway
    - 定期実行： Cloudwatch Events
    - データベース: DynamoDB
    - ストレージ: S3
    - 言語: golang
- CI/CD
    - GithubAction

## システム構成
Golang * Serverless Application Model * Line Messaging API

## 機能イメージ
### アプリ説明取得
### 睡眠時刻記録
### 起床時刻記録
### 5日分の睡眠記録取得
### 1ヶ月分の睡眠記録の取得
### 睡眠記録のグラフ表示