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

	messageIn, status, err := h.parseRequest(req)
	if err != nil {
		h.log.Error(err)
		http.Error(res, err.Error(), status)
		return
	}

	h.processRequest(*messageIn)
}

func (h *Handler) parseRequest(req *http.Request) (*interfaces.TelegramWebhookMessageStruct, int, error) {
	if req.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed: %s (actual) != POST (expected)", req.Method)
	}

	body := &interfaces.TelegramWebhookBodyStruct{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to decode request body: %s", err.Error())
	}

	if body.Message.Chat.ID != h.cfg.Telegram.ChatID {
		return nil, http.StatusOK, fmt.Errorf("chat id mismatch: %d (actual) != %d (expected)", body.Message.Chat.ID, h.cfg.Telegram.ChatID)
	}

	return &body.Message, 0, nil
}

func (h *Handler) processRequest(messageIn interfaces.TelegramWebhookMessageStruct) {
	h.matchers.Process(messageIn)
	h.state.ProcessMessage(messageIn)
}
