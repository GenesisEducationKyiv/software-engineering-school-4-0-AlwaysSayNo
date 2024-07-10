package cmd

import (
	"github.com/AlwaysSayNo/genesis-currency-api/common/pkg/envs"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/config"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	envs.Init("./pkg/common/envs/.env")

	//mailer := getMailer(config.LoadMailerConfig())

	waitServerWorking()
}

func getMailer(cnf config.MailerConfig) *service.Mailer {
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
