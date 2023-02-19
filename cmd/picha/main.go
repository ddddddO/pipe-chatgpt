package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	picha "github.com/ddddddO/pipe-chatgpt"
)

func main() {
	if err := do(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func do() error {
	apiKey := os.Getenv("CHATGPT_API_KEY")
	if apiKey == "" {
		if err := survey.AskOne(
			&survey.Password{
				Message: "ChatGPTのAPI keyを入力してください",
			}, &apiKey, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}
	gptClient := picha.NewGPTClient(apiKey)

	answerType := ""
	if err := survey.AskOne(
		&survey.Select{
			Message: "あなたの質問の種類は何ですか？",
			Options: []string{"テキストファイル", "音声", "テキスト"},
			Default: "テキスト",
		}, &answerType); err != nil {
		return err
	}

	responser, err := picha.ResponserFactory(answerType, gptClient)
	if err != nil {
		return err
	}
	return responser.Run()
}
