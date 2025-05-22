package types

type WebhookPayload struct {
	Type              string            `json:"type"`
	NewChatbotMessage NewChatbotMessage `json:"new_chatbot_message"`
}

type NewChatbotMessage struct {
	Text         string       `json:"text"`
	Conversation Conversation `json:"conversation"`
	Sender       Sender       `json:"sender"`
	Type         string       `json:"type"`
}

type Conversation struct {
	ID string `json:"id"`
}

type Sender struct {
	Type string `json:"type"`
}
