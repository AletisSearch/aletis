package message

import (
	"github.com/openai/openai-go/v3"
)

type MessageData struct {
	Age    int
	Gender string
	Year   int
}
type AiMessage []openai.ChatCompletionMessageParamUnion
type UserAssistant struct {
	User      string
	Assistant string
}

func (m *AiMessage) AddUserAssistant(us []UserAssistant) {
	for _, v := range us {
		(*m) = append(*m, openai.UserMessage(v.User), openai.AssistantMessage(v.Assistant))
	}
}
func (m *AiMessage) AddUser(userMsg string) {
	(*m) = append(*m, openai.UserMessage(userMsg))
}
