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

// mailerClient
//	- has an emailConsumer:
//		- subscribes to RabbitMQ and listens to events
//		- has listeners slice - functions which process messages from the queue
// 		- Subscribe - subscribes a new consumer to the listeners slice
// 		- Listen - function that listens either to queue or to the stop channel
// 		- handleMessage - sends messages to the listeners
//		- Close - closes connections to the queue
//	- Subscribe - receives a mailer, creates a function and passes it to the consumer
//	- Close - stops emailConsumer#Listen routine and calls emailConsumer#Close function
//	- defines an interface for mailer command

// + We save in local db next information: current currency (just db), users' emails;
// All information we receive from currency-rate service through commands from RabbitMQ;
// On currency-rate side we can either for each repository create a decorator, which will store in the target repository
// data and after that publish (this service might be responsible for SAGA)
// or this responsibility might take the service calling repositories;
// Remove current mail publishing command and its appropriate listener;
// Instead of them create 2 new commands and 2 listeners;
// + Each service should spin up its own db:(;
// + For email-service its own migrations should be created;
// + Create a repository for email-service
func main() {
	envs.Init("./pkg/envs/.env")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// DATABASE
	dbURL := db.GetDatabaseURL(dbconf.LoadDatabaseConfig())
	d := db.Init(dbURL)

	// MODULES
	mailModule := mailmodule.Init(d, emailconf.LoadEmailServiceConfig())

	mailTransport := getMailer(mailcnf.LoadMailerConfig())

	brokerClient, err := broker.NewClient(conscnf.LoadConsumerConfig())
	if err != nil {
		log.Fatalf("making mailer client: %v", err)
	}

	err = brokerClient.SubscribeCurrencyUpdateEvent(ctx, mailModule.CurrencyService)
	if err != nil {
		log.Fatalf("subscribing on CurrencyUpdateEvent: %v", err)
	}

	err = brokerClient.SubscribeUserSubscribedEvent(ctx, mailModule.UserService)
	if err != nil {
		log.Fatalf("subscribing on UserSubscribedEvent: %v", err)
	}

	err = brokerClient.SubscribeUserSubscriptionUpdatedEvent(ctx, mailModule.UserService)
	if err != nil {
		log.Fatalf("subscribing on UserSubscriptionUpdatedEvent: %v", err)
	}

	// ENGINE
	r := gin.Default()

	// HANDLERS
	registerRoutes(r, mailModule.EmailHandler)

	allJobs := scheduler.StartAllJobs(
		job.GetSendCurrencyEmailsJob(ctx, mailModule.EmailService, *mailTransport),
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

func registerRoutes(r *gin.Engine, emailHandler mailmodule.EmailHandler) {
	mailGroup := r.Group("/api/v1/mail")
	mailGroup.POST("/currency/send", emailHandler.SendCurrencyPriceEmails)
}

func startServer(r *gin.Engine) *http.Server {
	cnf := servcnf.LoadServerConfigConfig()
	server := &http.Server{
		Addr:    cnf.ApplicationPort,
		Handler: r.Handler(),
	}

	log.Println("Starting server")
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
