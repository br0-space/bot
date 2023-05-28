package interfaces

import telegramclient "github.com/br0-space/bot-telegramclient"

type MatcherRegistryInterface interface {
	Process(messageIn telegramclient.WebhookMessageStruct)
}

type MatcherInterface interface {
	IsEnabled() bool
	Identifier() string
	Help() []MatcherHelpStruct
	DoesMatch(messageIn telegramclient.WebhookMessageStruct) bool
	GetCommandMatch(messageIn telegramclient.WebhookMessageStruct) []string
	GetInlineMatches(messageIn telegramclient.WebhookMessageStruct) []string
	Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error)
	HandleError(messageIn telegramclient.WebhookMessageStruct, identifier string, err error)
}

type MatcherHelpStruct struct {
	Command     string
	Description string
	Usage       string
	Example     string
}
