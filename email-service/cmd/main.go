package main

import (
	"context"
	"errors"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/db"
	dbconf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/db/config"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/job"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker"
	conscnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker/consumer"
	mailmodule "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail"
	mailcnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/console"
	emailconf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/email/config"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/netsmtp"
	servcnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/server/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	envs.Init("./pkg/envs/.env")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// DATABASE
	dbURL := db.GetDatabaseURL(dbconf.LoadDatabaseConfig())
	d := db.Init(dbURL)

	// MODULES
	mailTransport := getMailer(mailcnf.LoadMailerConfig())

	mailModule := mailmodule.Init(d, mailTransport, emailconf.LoadEmailServiceConfig())

	brokerClient, err := broker.NewClient(conscnf.LoadConsumerConfig())
	if err != nil {
		log.Fatalf("making mailer client: %v", err)
	}

	err = brokerClient.SubscribeCurrencyUpdateEvent(ctx, &mailModule.CurrencyService)
	if err != nil {
		log.Fatalf("subscribing on CurrencyUpdateEvent: %v", err)
	}

	err = brokerClient.SubscribeUserSubscribedEvent(ctx, &mailModule.UserService)
	if err != nil {
		log.Fatalf("subscribing on UserSubscribedEvent: %v", err)
	}

	err = brokerClient.SubscribeUserSubscriptionUpdatedEvent(ctx, &mailModule.UserService)
	if err != nil {
		log.Fatalf("subscribing on UserSubscriptionUpdatedEvent: %v", err)
	}

	// ENGINE
	r := gin.Default()

	// HANDLERS
	registerRoutes(r, &mailModule.EmailHandler)

	allJobs := scheduler.StartAllJobs(
		job.GetSendCurrencyEmailsJob(ctx, &mailModule.EmailService, *mailTransport),
	)

	server := startServer(r)
	waitServerWorking()

	// STOP SERVER
	gracefulShutdown(ctx, allJobs, server, *brokerClient)
}

func getMailer(cnf mailcnf.MailerConfig) *service.Mailer {
	var mailer service.Mailer

	switch cnf.Type {
	case "console":
		mailer = console.NewConsoleMailer(cnf)
	case "net/smtp":
		mailer = netsmtp.NewNetSMTPMailer(cnf)
	default:
		mailer = console.NewConsoleMailer(cnf)
	}

	return &mailer
}

func registerRoutes(r *gin.Engine, emailHandler mailmodule.Handler) {
	mailGroup := r.Group("/api/v1/mail")
	mailGroup.POST("/currency/send", emailHandler.SendCurrencyPriceEmails)
}

func startServer(r *gin.Engine) *http.Server {
	cnf := servcnf.LoadServerConfigConfig()
	server := &http.Server{
		Addr:    cnf.ApplicationPort,
		Handler: r.Handler(),
	}

	log.Println("Starting server on port:", cnf.ApplicationPort)
	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server is started")

	return server
}

func waitServerWorking() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Println("Server received next signal:", sign.String())
}

func gracefulShutdown(ctx context.Context, allJobs *cron.Cron, server *http.Server, brokerClient broker.Client) {
	log.Println("Stopping server")

	cnf := servcnf.LoadServerConfigConfig()
	waitSeconds := cnf.GracefulShutdownWaitTimeSeconds

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Duration(waitSeconds)*time.Second)
	defer shutdownCancel()

	// STOP JOBS
	allJobs.Stop()

	// STOP WEB SERVER
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}

	// STOP MAILER CLIENT
	if err := brokerClient.Close(); err != nil {
		log.Printf("mailer shutdown: %v", err)
	}

	select {
	case <-shutdownCtx.Done():
		log.Printf("Timeout of %d seconds\n", waitSeconds)
	}

	log.Println("Server exiting")
}
