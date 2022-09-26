package telegram

import (
	"context"
	"fmt"

	"github.com/defrell01/telebot/pkg/repository"
)

func (b *Bot) generateAuthorizationLink(chatID uint64) (string, error) {
	redirectURL := b.generateRedirectURL(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectURL(chatID uint64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
