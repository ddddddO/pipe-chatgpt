package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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
		questionAPIKey := &survey.Password{
			Message: "ChatGPTのAPI keyを入力してください",
		}
		if err := survey.AskOne(questionAPIKey, &apiKey, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}
	gptClient := picha.NewGPTClient(apiKey)

	questionType := &survey.Select{
		Message: "あなたの質問の種類は何ですか？",
		Options: []string{"テキストファイル", "音声", "テキスト"},
		Default: "テキスト",
	}
	answerType := ""
	if err := survey.AskOne(questionType, &answerType); err != nil {
		return err
	}

	// TODO: ここは、なにかinterfaceで定義して、その実装を呼び出す、に変更したいかも
	switch answerType {
	case "テキスト":
		for {
			questionText := &survey.Input{
				Message: "聞きたいことは？",
			}
			answerText := ""
			if err := survey.AskOne(questionText, &answerText, survey.WithValidator(survey.Required)); err != nil {
				return err
			}
			if err := gptClient.RequestToDavinci(answerText); err != nil {
				return err
			}
		}
	case "テキストファイル":
		questionPath := &survey.Input{
			Message: "そのファイルパスを入力してください",
		}
		answerPath := ""
		if err := survey.AskOne(questionPath, &answerPath, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
		path, err := filepath.Abs(answerPath)
		if err != nil {
			return err
		}
		if !isExist(path) {
			return fmt.Errorf("no exist file: %s", answerPath)
		}

		questionProcessingFile := &survey.Input{
			Message: "このファイルをどうしたいですか？",
		}
		answerProcessingFile := ""
		if err := survey.AskOne(questionProcessingFile, &answerProcessingFile, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// TODO: 一旦、jsonとかcsvとか考えずに普通のテキストファイルとして作る
		sc := bufio.NewScanner(f)
		in := "「"
		for sc.Scan() {
			in += sc.Text()
		}
		if err := sc.Err(); err != nil {
			return err
		}
		in += "」"
		in += answerProcessingFile

		return gptClient.RequestToDavinci(in)
	case "音声":
		// TODO:
	default:
		return fmt.Errorf("未定義の種類です")
	}

	return nil
}

func isExist(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
