# sleep-recorder
<img width="290" alt="スクリーンショット 2021-11-01 21 29 42" src="https://user-images.githubusercontent.com/39262724/139671550-6f9a7100-757a-4ae6-a6ca-d59349c26e11.png">
このアプリでは，LINEを通して簡単に睡眠時間をお手軽に記録し、1週間・1ヶ月の**睡眠情報**やその**評価**を**簡単に可視化**する

- 対象: 自分の睡眠状態を簡単に記録・評価・可視化したい人向けのアプリケーション
- 目的: このアプリを通して自身の睡眠を振り返り、規則正しく、理想とする睡眠習慣・生活習慣を掴み取ること

## バッジ
[![CI](https://github.com/KazuwoKiwame12/sleep-recorder/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/KazuwoKiwame12/sleep-recorder/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/KazuwoKiwame12/sleep-recorder)](/LICENSE)

## リリース
Coming soon...  
issueとかで機能の要望などを送ってもらえると嬉しいです！  
また、PRなどもいただけると嬉しいです

## 主機能
1. ユーザがLINE画面で"眠る"を打つと、システムは睡眠時刻を保存する
2. ユーザがLINE画面で"起きた"を打つと、システムは起床時間と1との間の睡眠時間を記録する
3. ユーザがLINE画面で"取得"を打つと、システムは現時点から4日前までの5日間の睡眠記録リストを表示する
4. ユーザがLINE画面で"取得 年(数字) 月(数字)"を打つ、システムは該当する月の睡眠記録リストを表示する
6. システムは1週間に1度、月曜日の夜23時に、LINEにその週の睡眠記録リストのグラフを送信する

### 評価
睡眠記録の評価に関しては4段階存在しています。
- 3ポイント: "🤩 完璧!"
- 2ポイント: "😁 良いね!"
- 1ポイント: "😥 がんばれ!"
- 0ポイント: "😱 伸び代しかない!"   
数字が大きいほど評価が高いです。  
数字が小さく評価が悪くとも、ネガティブにならないような言葉が添えてあります。

評価の仕方について、最初は3ポイント所持しており、
以下の状態になると、その都度1ポイント減ります。
1. 1時以降の睡眠の場合
2. 起床時間が8時以上
3. 睡眠時間が7.0時間~7.9時間ではない

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
**Golang * Serverless Application Model * Line Messaging API**

![アーキテクチャ](https://user-images.githubusercontent.com/39262724/139680118-66d677bc-8ba5-4780-b1ca-aa123dcb1841.png)


## 機能イメージ
### アプリ説明取得
<img width="1005" alt="explain_1" src="https://user-images.githubusercontent.com/39262724/139667477-a04da8ff-aaaa-4352-a80c-2ded06eaa75c.png">
<img width="1013" alt="explain_2" src="https://user-images.githubusercontent.com/39262724/139667527-82887f21-49fe-4573-9b4d-fc54028b05e5.png">
<img width="1006" alt="explain_3" src="https://user-images.githubusercontent.com/39262724/139667563-3f99ed69-81d0-42cb-a5fe-5756505664e5.png">

### 睡眠時刻記録
<img width="1005" alt="睡眠記録_1" src="https://user-images.githubusercontent.com/39262724/139667631-0350f29e-2765-4d62-8c1c-849ec2629132.png">
<img width="1008" alt="睡眠記録_2" src="https://user-images.githubusercontent.com/39262724/139667642-1e080410-fa91-4779-aee3-54d7b1b5a9e7.png">

※ 22時から"睡眠時刻"の記録が可能になります 

※ 入力の度に"睡眠時刻"の記録が上書きされます

### 起床時刻記録
<img width="1009" alt="起床記録_1" src="https://user-images.githubusercontent.com/39262724/139667717-a930e632-8e77-41f8-8d22-a5c651e9b90a.png">
<img width="1005" alt="起床記録_2" src="https://user-images.githubusercontent.com/39262724/139667723-f8f322be-18c4-475d-a913-0970c95aec44.png">


### 5日分の睡眠記録取得
<img width="1005" alt="取得_1" src="https://user-images.githubusercontent.com/39262724/139667812-d5f1f7c7-d372-4011-85e3-02c433288658.png">
<img width="1002" alt="取得_2" src="https://user-images.githubusercontent.com/39262724/139667821-c1120f82-4012-4933-a3d6-8bd65bba7a09.png">

※ 期間内にデータがない場合は取得できません

### 1ヶ月分の睡眠記録の取得
<img width="1000" alt="取得_月_1" src="https://user-images.githubusercontent.com/39262724/139667860-afcaa9ee-3460-4a3a-90ca-9b5ffcca7399.png">
<img width="1003" alt="取得_月_2" src="https://user-images.githubusercontent.com/39262724/139667912-41d41780-95b2-4160-aba0-562a14ccb6ac.png">

※ このテキストは全角でも半角でも問題なく動作します 

※ 期間内にデータがない場合は取得できません

### 睡眠記録のグラフ表示
<img width="1004" alt="グラフ_1" src="https://user-images.githubusercontent.com/39262724/139667955-40eac79c-7fc2-4bca-9309-238a9af93ce2.png">
<img width="1004" alt="グラフ_2" src="https://user-images.githubusercontent.com/39262724/139667963-bf1e22f8-2bd1-4ef1-b8c1-8f469bf6b518.jpg">


※ 現時点(2021/11/1)のグラフのUIがあまり良くないと考えているので、今後改善していきます

※ データがない日の部分は表示されません
