package utility

import (
	"regexp"
	"strings"
)

const (
	MessageHelp string = `このアプリには主に3つの機能が存在します。
"1": 睡眠記録の登録
"2": 睡眠記録の取得
"3": 睡眠記録の定期報告
	
それぞれについての概要と使い方の説明を記述していきます。

"1"は睡眠時刻と起床時刻を記録できます。
【以下の文字の送信】
1. "眠る"
	睡眠時刻の記録
2. "起きた"
	起床時刻の記録
	睡眠時間の記録

"2"は5日分の睡眠記録と、指定した月の睡眠記録を取得できます。
【以下の文字の送信】
1. "取得"
	5日分(本日~4日前)の睡眠記録の取得
2. "取得 (年) (月)"
	指定した年月の1ヶ月の睡眠記録の取得
	()の中には数字

"3"は、毎週月曜日の夜9時に1週間分の睡眠記録の遷移グラフの画像が送信されます。
	`
	MessageDefault string = `アプリの使い方は"説明"を送信することで得られます。

【使い方概要】
1. "眠る" 
	睡眠時刻を記録

2. "起きた"
	起床時刻を記録
	睡眠時間の記録

3. "取得"
	5日分の睡眠記録の取得

4. "取得 (年) (月)"
	指定した年月の1ヶ月の睡眠記録の取得
	`
	MessageSuccessRecord  string = "記録できました"
	MessageSystemError    string = "システム側でエラーが発生しました"
	MessageBedinTimeError string = "12時から22時の間に眠ることはありません。"
	MessageNotFound       string = "記録が存在しません"
	MessageNotSleep       string = "本日は眠っていないので起床時刻の記録ができません"
)

type Command int

const (
	CommandBedin Command = iota
	CommandWake
	CommandFiveDays
	CommandMonth
	CommandHelp
	CommandDefault
)

func ValidateCommand(com string) Command {
	if com == "眠る" {
		return CommandBedin
	}
	if com == "起きた" {
		return CommandWake
	}
	if com == "説明" {
		return CommandHelp
	}
	if regexp.MustCompile(`取得.*`).MatchString(com) {
		slice := strings.Split(com, " ")
		if len(slice) == 1 {
			return CommandFiveDays
		}
		return CommandMonth
	}
	return CommandDefault
}
