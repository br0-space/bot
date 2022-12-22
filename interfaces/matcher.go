package interfaces

type MatcherRegistryInterface interface {
	Init()
	Process(messageIn TelegramWebhookMessageStruct)
}

type MatcherInterface interface {
	IsEnabled() bool
	GetIdentifier() string
	GetHelp() []MatcherHelpStruct
	DoesMatch(messageIn TelegramWebhookMessageStruct) bool
	GetCommandMatch(messageIn TelegramWebhookMessageStruct) []string
	GetInlineMatches(messageIn TelegramWebhookMessageStruct) []string
	Process(messageIn TelegramWebhookMessageStruct) (*[]TelegramMessageStruct, error)
	HandleError(messageIn TelegramWebhookMessageStruct, identifier string, err error)
}

type MatcherHelpStruct struct {
	Command     string
	Description string
	Usage       string
	Example     string
}
