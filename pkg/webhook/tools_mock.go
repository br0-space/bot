package webhook

type MockTools struct{}

func NewMockTools() *MockTools {
	return &MockTools{}
}

func (t *MockTools) SetWebhookURL() {}
