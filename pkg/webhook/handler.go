package webhook

import (
	"encoding/json"
	"fmt"
	logger "github.com/br0-space/bot-logger"
	"github.com/br0-space/bot/interfaces"
	"net/http"
)

type Handler struct {
	log      logger.Interface
	cfg      *interfaces.ConfigStruct
	matchers interfaces.MatcherRegistryInterface
	state    interfaces.StateServiceInterface
}

func NewHandler(
	config *interfaces.ConfigStruct,
	matchers interfaces.MatcherRegistryInterface,
	state interfaces.StateServiceInterface,
) *Handler {
	return &Handler{
		log:      logger.New(),
		cfg:      config,
		matchers: matchers,
		state:    state,
	}
}

func (h *Handler) InitMatchers() {
	h.matchers.Init()
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	h.log.Debugf("%s %s %s from %s", req.Method, req.URL, req.Proto, req.RemoteAddr)

	messageIn, err, status := h.parseRequest(req)
	if err != nil {
		h.log.Error(err)
		http.Error(res, err.Error(), status)
		return
	}

	h.processRequest(*messageIn)
}

func (h *Handler) parseRequest(req *http.Request) (*interfaces.TelegramWebhookMessageStruct, error, int) {
	if req.Method != http.MethodPost {
		return nil, fmt.Errorf("method not allowed: %s (actual) != POST (expected)", req.Method), http.StatusMethodNotAllowed
	}

	body := &interfaces.TelegramWebhookBodyStruct{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return nil, fmt.Errorf("unable to decode request body: %s", err.Error()), http.StatusBadRequest
	}

	if body.Message.Chat.ID != h.cfg.Telegram.ChatID {
		return nil, fmt.Errorf("chat id mismatch: %d (actual) != %d (expected)", body.Message.Chat.ID, h.cfg.Telegram.ChatID), http.StatusOK
	}

	return &body.Message, nil, 0
}

func (h *Handler) processRequest(messageIn interfaces.TelegramWebhookMessageStruct) {
	h.matchers.Process(messageIn)
	h.state.ProcessMessage(messageIn)
}
