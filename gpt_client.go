package picha

import (
	"context"
	"fmt"

	"github.com/PullRequestInc/go-gpt3"
)

type gptClient struct {
	client gpt3.Client
}

func NewGPTClient(apiKey string) *gptClient {
	return &gptClient{
		client: gpt3.NewClient(apiKey),
	}
}

// 各エンジンで得意なことが違うそう: https://qiita.com/kinuta_masa/items/ca6ab79bfca438f8f970
// TODO: そのため、ユーザーのリクエストによってエンジンを変えたい。
func (g *gptClient) RequestToAda(request string) error {
	return g.requestTo(gpt3.TextAda001Engine, request)
}

func (g *gptClient) RequestToBabbage(request string) error {
	return g.requestTo(gpt3.TextBabbage001Engine, request)
}

func (g *gptClient) RequestToCurie(request string) error {
	return g.requestTo(gpt3.TextCurie001Engine, request)
}

func (g *gptClient) RequestToDavinci(request string) error {
	return g.requestTo(gpt3.TextDavinci003Engine, request)
}

func (g *gptClient) requestTo(who, request string) error {
	err := g.client.CompletionStreamWithEngine(
		context.Background(),
		who,
		gpt3.CompletionRequest{
			Prompt: []string{
				request,
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
	return nil
}
