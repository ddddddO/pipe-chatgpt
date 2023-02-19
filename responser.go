package picha

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

type Responser interface {
	Run() error
}

func ResponserFactory(answerType string, gptClient *gptClient) (Responser, error) {
	switch answerType {
	case "テキスト":
		return &textResponser{gptClient: gptClient}, nil
	case "テキストファイル":
		return &textFileResponser{gptClient: gptClient}, nil
	case "音声":
		return &voiceResponser{gptClient: gptClient}, nil
	}
	return nil, fmt.Errorf("未定義の種類です")
}

type textResponser struct {
	gptClient *gptClient
}

func (tr *textResponser) Run() error {
	for {
		answerText := ""
		if err := survey.AskOne(
			&survey.Input{
				Message: "聞きたいことは？",
			}, &answerText, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
		if err := tr.gptClient.RequestToDavinci(answerText); err != nil {
			return err
		}
	}
	return nil
}

type textFileResponser struct {
	gptClient *gptClient
}

func (tfr *textFileResponser) Run() error {
	answerPath := ""
	if err := survey.AskOne(
		&survey.Input{
			Message: "そのファイルパスを入力してください",
		}, &answerPath, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	path, err := filepath.Abs(answerPath)
	if err != nil {
		return err
	}
	if !isExist(path) {
		return fmt.Errorf("no exist file: %s", answerPath)
	}

	answerProcessingFile := ""
	if err := survey.AskOne(
		&survey.Input{
			Message: "このファイルをどうしたいですか？",
		}, &answerProcessingFile, survey.WithValidator(survey.Required)); err != nil {
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

	return tfr.gptClient.RequestToDavinci(in)
}

func isExist(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

type voiceResponser struct {
	gptClient *gptClient
}

func (vr *voiceResponser) Run() error {
	// TODO:
	return nil
}
