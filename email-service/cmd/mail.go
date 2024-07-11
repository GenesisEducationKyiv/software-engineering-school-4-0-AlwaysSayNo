package cmd

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/common/pkg/envs"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker/client"
	conscnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker/consumer/config"
	mailcnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	servcnf "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/server/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	envs.Init("./pkg/common/envs/.env")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mailer := getMailer(mailcnf.LoadMailerConfig())

	mailerClient, err := client.NewClient(conscnf.LoadConsumerConfig())
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
		mailer = service.NewConsoleMailer(cnf)
	case "net/smtp":
		mailer = service.NewNetSMTPMailer(cnf)
	}

	return &mailer
}

func waitServerWorking() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Println("Server received next signal:", sign.String())
}

func gracefulShutdown(ctx context.Context, mailerClient client.Client) {
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
