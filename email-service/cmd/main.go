package main

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker"
	conscnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker/consumer/config"
	mailcnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/console"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/netsmtp"
	servcnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/server/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"
	"log"
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

func main() {
	envs.Init("./pkg/envs/.env")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mailer := getMailer(mailcnf.LoadMailerConfig())

	mailerClient, err := broker.NewClient(conscnf.LoadConsumerConfig())
	if err != nil {
		log.Fatalf("making mailer client: %v", err)
	}

	err = mailerClient.Subscribe(ctx, *mailer)
	if err != nil {
		log.Fatalf("subscribing mailer: %v", err)
	}

	waitServerWorking()

	// STOP SERVER
	gracefulShutdown(ctx, *mailerClient)
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

func waitServerWorking() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Println("Server received next signal:", sign.String())
}

func gracefulShutdown(ctx context.Context, mailerClient broker.Client) {
	log.Println("Stopping server")

	cnf := servcnf.LoadServerConfigConfig()
	waitSeconds := cnf.GracefulShutdownWaitTimeSeconds

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Duration(waitSeconds)*time.Second)
	defer shutdownCancel()

	// STOP MAILER CLIENT
	if err := mailerClient.Close(); err != nil {
		log.Printf("mailer shutdown: %v", err)
	}

	select {
	case <-shutdownCtx.Done():
		log.Printf("Timeout of %d seconds\n", waitSeconds)
	}

	log.Println("Server exiting")
}
