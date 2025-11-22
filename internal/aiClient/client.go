//go:generate msgp -tests=false

package aiclient

import (
	"context"
	"encoding/json/v2"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/AletisSearch/aletis/internal/cache"
	"github.com/AletisSearch/aletis/internal/db"
	"github.com/AletisSearch/aletis/internal/message"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Client struct {
	apiClient openai.Client
	cache     *cache.Cache[Output, *Output]
}
type CompletionUsage struct {
	Cost        float64 `json:"cost"`
	IsByoK      bool    `json:"is_byok"`
	CostDetails struct {
		UpstreamInferenceCost            float64 `json:"upstream_inference_cost"`
		UpstreamInferencePromptCost      float64 `json:"upstream_inference_prompt_cost"`
		UpstreamInferenceCompletionsCost float64 `json:"upstream_inference_completions_cost"`
	} `json:"cost_details"`
}

type Output struct {
	Content string  `json:"content"`
	Cost    float64 `json:"cost"`
}

func NewClient(baseurl, openaiKey string, db *db.Queries) *Client {
	c := &Client{
		apiClient: openai.NewClient(
			option.WithBaseURL(baseurl),
			option.WithAPIKey(openaiKey),
			//option.WithHeader()
			option.WithJSONSet("usage.include", true),
			option.WithHeader("HTTP-Referer", "https://github.com/AletisSearch/aletis"),
			option.WithHeader("X-Title", "Aletis"),
		),
		cache: cache.New[Output](db),
	}
	return c
}

func (c *Client) RunQueryExpand(ctx context.Context, q string) (*Output, error) {
	var sysMsg strings.Builder
	mData := message.MessageData{
		Year: time.Now().Year(),
		// Age:    22,
		// Gender: "Male",
	}
	sysMsg.WriteString(message.SystemQueryExpand(3, 5, mData))
	cacheKey := "aiClient-" + q

	o, err := c.cache.Get(ctx, cacheKey)
	if err == nil {
		slog.Info("Cache Hit", "Key", cacheKey)
		return o, nil
	}
	if !errors.Is(err, cache.ErrNotFoundInCache) && !errors.Is(err, cache.ErrOldCache) {
		return nil, err
	}

	us, err := message.TemplateToUserAssistant(message.QueryExpandData, mData)
	if err != nil {
		return &Output{}, err
	}
	out, err := c.Run(ctx, "google/gemma-3-12b-it", sysMsg.String(), q, us...)
	if err != nil {
		return nil, err
	}
	return out, c.cache.Set(ctx, cacheKey, out, time.Hour*25)
}

func (c *Client) Run(ctx context.Context, model, system, query string, messages ...message.UserAssistant) (*Output, error) {
	m := message.AiMessage{
		openai.SystemMessage(system),
	}
	if len(messages) > 0 {
		m.AddUserAssistant(messages)
	}
	m.AddUser(query)
	chatCompletion, err := c.apiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: m,
		//Model:    "google/gemini-2.5-flash-lite-preview-09-2025",
		Model: model,
		//ReasoningEffort: "minimal",
	})
	if err != nil {
		return nil, err
	}

	cost := CompletionUsage{}
	if err = json.Unmarshal([]byte(chatCompletion.Usage.RawJSON()), &cost); err != nil {
		return nil, err
	}

	return &Output{Content: chatCompletion.Choices[0].Message.Content, Cost: cost.Cost}, nil
}
