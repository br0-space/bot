package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"net/http"
)

type Handler struct {
	Log      interfaces.LoggerInterface
	Cfg      *interfaces.ConfigStruct
	Matchers interfaces.MatcherRegistryInterface
}

func NewHandler(logger interfaces.LoggerInterface, config *interfaces.ConfigStruct, matchers interfaces.MatcherRegistryInterface) *Handler {
	return &Handler{
		Log:      logger,
		Cfg:      config,
		Matchers: matchers,
	}
}

func (h *Handler) InitMatchers() {
	h.Matchers.Init()
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	h.Log.Debugf("%s %s %s from %s", req.Method, req.URL, req.Proto, req.RemoteAddr)

	switch req.Method {
	case "POST":
	default:
		h.Log.Error("method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	messageIn, err := h.ParseRequest(req)
	if err != nil {
		h.Log.Error(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if messageIn.Chat.ID != h.Cfg.Telegram.ChatID {
		h.Log.Warningf("chat id mismatch: %d (actual) != %d (expected)", messageIn.Chat.ID, h.Cfg.Telegram.ChatID)
		http.Error(res, "chat id mismatch", http.StatusOK)
		return
	}

	h.Matchers.Process(*messageIn)
}

func (h *Handler) ParseRequest(req *http.Request) (*interfaces.TelegramWebhookMessageStruct, error) {
	body := &interfaces.TelegramWebhookBodyStruct{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return nil, fmt.Errorf("unable to decode request body: %s", err.Error())
	}

	return &body.Message, nil
}
