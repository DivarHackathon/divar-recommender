package services

import (
	"bytes"
	"divar_recommender/internal/types"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatService struct {
	apiKey  string
	baseURL string
}

func NewChatService(apiKey string) *ChatService {
	return &ChatService{
		apiKey:  apiKey,
		baseURL: "https://open-api.divar.ir/experimental/open-platform/chat/bot/conversations",
	}
}

func (s *ChatService) BuildImagePreview(ad types.Ad) types.ChatMessage {
	return types.ChatMessage{
		Type:        "TEXT",
		TextMessage: " ",
		Buttons: types.ButtonsWrapper{
			Rows: []types.ButtonRow{
				{
					Buttons: []types.Button{
						{
							Action: types.Action{
								OpenDirectLink: ad.Image,
							},
							IconName: "CAR",
							Caption:  "📸 مشاهده تصویر",
						},
					},
				},
			},
		},
	}
}

func (s *ChatService) BuildTextOnly(ad types.Ad) types.ChatMessage {
	return types.ChatMessage{
		Type:        "TEXT",
		TextMessage: fmt.Sprintf("📌 %s\n💰 قیمت: %d تومان", ad.Title, ad.Price),
	}
}

func (s *ChatService) BuildLinkButton(ad types.Ad) types.ChatMessage {
	return types.ChatMessage{
		Type:        "TEXT",
		TextMessage: "📎 مشاهده آگهی در دیوار",
		Buttons: types.ButtonsWrapper{
			Rows: []types.ButtonRow{
				{
					Buttons: []types.Button{
						{
							Action:   types.Action{OpenDirectLink: fmt.Sprintf("https://divar.ir/v/%s", ad.Token)},
							IconName: "CAR",
							Caption:  "📲 باز کردن آگهی",
						},
					},
				},
			},
		},
	}
}

func (s *ChatService) SendMessage(conversationID string, message types.ChatMessage) error {
	url := fmt.Sprintf("%s/%s/messages", s.baseURL, conversationID)

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *ChatService) SendAdDetails(conversationID string, ad types.Ad) error {
	textMsg := s.BuildTextOnly(ad)
	if err := s.SendMessage(conversationID, textMsg); err != nil {
		return err
	}

	if ad.Image != "" {
		imageMsg := s.BuildImagePreview(ad)
		if err := s.SendMessage(conversationID, imageMsg); err != nil {
			return err
		}
	}

	linkMsg := s.BuildLinkButton(ad)
	if err := s.SendMessage(conversationID, linkMsg); err != nil {
		return err
	}

	return nil
}
