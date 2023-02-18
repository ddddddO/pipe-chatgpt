package main

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/PullRequestInc/go-gpt3"
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
	gptClient := gpt3.NewClient(apiKey)

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

			err := gptClient.CompletionStreamWithEngine(
				context.Background(),
				gpt3.TextDavinci003Engine,
				gpt3.CompletionRequest{
					Prompt: []string{
						answerText,
					},
					MaxTokens:   gpt3.IntPtr(3000),
					Temperature: gpt3.Float32Ptr(0),
				}, func(resp *gpt3.CompletionResponse) {
					fmt.Print(resp.Choices[0].Text)
				})
			if err != nil {
				return err
			}
			fmt.Println()
		}
	case "テキストファイル":
		questionPath := &survey.Input{
			Message: "そのファイルパスを入力してください",
		}
		answerPath := ""
		if err := survey.AskOne(questionPath, &answerPath, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
		// TODO:
		fmt.Printf("あなたが選んだファイル: %s\n", answerPath)
	case "音声":
		// TODO:
	default:
		return fmt.Errorf("未定義の種類です")
	}

	return nil
}
