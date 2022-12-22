package main

import (
	"github.com/br0-space/bot/container"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"net/http"
	"time"
)

func main() {
	pflag.Parse()

	logger := container.ProvideLogger()
	config := container.ProvideConfig()
	databaseMigration := container.ProvideDatabaseMigration()
	webhookHandler := container.ProvideTelegramWebhookHandler()
	webhookTools := container.ProvideTelegramWebhookTools()

	//logger.Info("Initializing Prometheus stats collector")
	//
	//prometheusHandler := prometheus.Handler{}
	//stats.Register(&prometheusHandler)
	//defer stats.Flush()
	//
	//// Start a new collector for the current process, reporting Go metrics.
	//procStatsCollector := procstats.StartCollector(procstats.NewGoMetrics())
	//// Gracefully stops stats collection.
	//defer procStatsCollector.Close()

	if config.Database.AutoMigrate {
		logger.Info("Running database migrations")
		if err := databaseMigration.Migrate(); err != nil {
			logger.Fatal(err)
		}
	}

	logger.Info("Initializing Matchers")
	webhookHandler.InitMatchers()

	logger.Info("Initializing Telegram webhook")
	webhookTools.SetWebhookURL()

	logger.Info("Starting HTTP server listening on", config.Server.ListenAddr)

	r := mux.NewRouter()
	r.HandleFunc("/webhook", webhookHandler.ServeHTTP)
	//r.HandleFunc("/metrics", prometheusHandler.ServeHTTP)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/", r)

	srv := &http.Server{
		Addr:           config.Server.ListenAddr,
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    0,
		MaxHeaderBytes: 4096,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func notFound(res http.ResponseWriter, req *http.Request) {
	logger := container.ProvideLogger()

	logger.Debugf("%s %s %s from %s", req.Method, req.URL, req.Proto, req.RemoteAddr)
	logger.Error("not found")
}
