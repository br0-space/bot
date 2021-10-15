package cmd

import (
	"net/http"
	"time"

	"github.com/br0-space/bot/container"
	"github.com/br0-space/bot/internal/config"
	"github.com/br0-space/bot/internal/logger"
	"github.com/br0-space/bot/internal/matcher/ping"
	"github.com/br0-space/bot/internal/matcher/registry"
	"github.com/br0-space/bot/internal/telegram/webhook"
	"github.com/gorilla/mux"
	"github.com/segmentio/stats/v4"
	"github.com/segmentio/stats/v4/procstats"
	"github.com/segmentio/stats/v4/prometheus"
)

type Cmd struct {
	logger logger.Interface
	cfg    config.Server
}

func NewCmd() *Cmd {
	return &Cmd{
		logger: container.ProvideLoggerService(),
		cfg:    container.ProvideConfig().Server,
	}
}

// Create an HTTP server listening for webhook requests from telegram on port 3000
func (c *Cmd) RunServer() error {
	// telegram.SetWebhookURL()

	prometheusHandler := prometheus.Handler{}
	stats.Register(&prometheusHandler)
	defer stats.Flush()

	// Start a new collector for the current process, reporting Go metrics.
	procstatsCollector := procstats.StartCollector(procstats.NewGoMetrics())
	// Gracefully stops stats collection.
	defer procstatsCollector.Close()

	r := mux.NewRouter()
	r.HandleFunc("/webhook", c.WebhookHandler).Methods("POST")
	r.HandleFunc("/metrics", prometheusHandler.ServeHTTP)
	http.Handle("/", r)

	// Handler:        httpstats.NewHandler(r),

	stats.Incr("test")

	srv := &http.Server{
		Addr:           c.cfg.ListenAddr,
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    0,
		MaxHeaderBytes: 4096,
	}

	c.logger.Info("listening on", srv.Addr)

	return srv.ListenAndServe()
}

// Handle a webhook request sent by telegram
func (c *Cmd) WebhookHandler(res http.ResponseWriter, req *http.Request) {
	// Parse the request
	messageIn, err := webhook.ParseRequest(res, req)
	if err != nil {
		c.logger.Error(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Send request message to all matchers in registry
	registry.NewRegistry().ProcessWebhookMessageInAllMatchers(*messageIn)
}

// Process a test message only (for development use)
func TestOnly() {
	// requestMessage := telegram.RequestMessage{}
	// requestMessage.From.ID = 666
	// requestMessage.From.Username = "Testuser"
	// requestMessage.From.FirstName = "Test"
	// requestMessage.From.LastName = "User"
	// requestMessage.Text = "weggesynced"
	// requestMessage.Chat.ID = 20551552
	//
	// matcher.ExecuteMatchers(requestMessage)
	//
	// logger.Log.Info("--- done ---")

	message := webhook.Message{Text: "foobar blafasel"}
	matcher := ping.MakeMatcher()
	matcher.ProcessMessage(message)
}
