package utility

import (
	"regexp"
	"strings"
)

const (
	MessageHelp string = `
	このアプリには3つの機能が存在します。\n
	1: 睡眠記録の登録\n
	2: 睡眠記録の取得\n
	3: 睡眠記録の定期報告\n\n
	それぞれについての概要と使い方の説明を記述していきます。\n
	まず"1"についてです。\n
	1は睡眠時刻と起床時刻を記録できます。\n
	使い方...以下の文字の送信\n
	\t睡眠時刻を記録: "眠る"\n
	\t起床時刻を記録: "起きた"\n
	次に"2""についてです。\n
	2は本日から4日前までの5日分の睡眠記録と、指定した月の睡眠記録を取得できます。\n
	使い方...以下の文字の送信\n
	\t5日分の睡眠記録の取得: "取得"\n
	\t指定した月の睡眠記録の取得: "取得 (年) (月)"...()の中には数字\n
	最後に"3"は、毎週月曜日の夜9時に1週間分の睡眠記録の遷移グラフの画像が送信されます。
	`
	MessageDefault string = `
	アプリの使い方は"説明"を送信することで得られます。\n
	使い方概要...以下の文字を送信すると以下の機能が実行される\n
	\t睡眠時刻を記録: "眠る"\n
	\t起床時刻を記録: "起きた"\n
	\t5日分の睡眠記録の取得: "取得"\n
	\t指定した月の睡眠記録の取得: "取得 (年) (月)"
	`
	MessageSuccessRecord  string = "記録できました"
	MessageSystemError    string = "システム側でエラーが発生しました"
	MessageBedinTimeError string = "12時から22時の間に眠ることはありません。"
	MessageNotFound       string = "記録が存在しません"
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
