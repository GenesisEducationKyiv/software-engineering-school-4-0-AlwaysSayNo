package main

import (
	"context"
	"errors"
	"github.com/AlwaysSayNo/genesis-currency-api/common/pkg/envs"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail"
	prodcnf "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail/producer/config"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/notifier"
	emailconf "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/notifier/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/db"
	dbconf "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/db/config"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/job"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/middleware"
	currencymodule "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency"
	ratecdnjsdelivr "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/external/rater/cdnjsdelivr"
	rategovua "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/external/rater/gov_ua"
	rateprivate "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/external/rater/private"
	currencyhand "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/handler"
	currencyconf "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/config"
	currencyserv "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/service"
	usermodule "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user"
	userhand "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/api/handler"
	userconf "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/server/config"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	envs.Init("./pkg/common/envs/.env")

	// CONTEXT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// DATABASE
	dbURL := db.GetDatabaseURL(dbconf.LoadDatabaseConfig())
	d := db.Init(dbURL)

	// CURRENCY API HANDLERS
	currencyProvider := getCurrencyProviderChain()

	// MODULES
	userModule := usermodule.Init(d)
	currencyModule := currencymodule.Init(currencyProvider)
	//emailModule := emailmodule.Init(userModule.Service, currencyModule.Service, emailconf.LoadEmailServiceConfig())

	// EMAIL PRODUCER CLIENT
	mailClient := getMailClient()
	mailNotifier := notifier.NewEmailNotifier(mailClient, currencyModule.Service, userModule.Service, emailconf.LoadEmailServiceConfig())

	// ENGINE
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	// JOBS
	scheduler := job.StartAllJobs(
		job.GetUpdateCurrencyJob(ctx, currencyModule.Service),
		job.GetSendEmailsJob(ctx, mailNotifier),
	)

	// HANDLERS
	currencyhand.RegisterRoutes(r, currencyModule.Handler)
	userhand.RegisterRoutes(r, userModule.Handler)
	//emailhand.RegisterRoutes(r, emailModule.Handler)

	// START SERVER
	server := startServer(r)
	waitServerWorking()

	// STOP SERVER
	gracefulShutdown(ctx, scheduler, server, mailClient)
}

func getCurrencyProviderChain() currencyserv.Provider {
	// GET PROVIDERS
	privateClient, err := rateprivate.NewClient(currencyconf.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Private Bank currency provider")
	}
	govUaClient, err := rategovua.NewClient(currencyconf.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Bank Gov Ua currency provider")
	}
	jsDelivrClient, err := ratecdnjsdelivr.NewClient(currencyconf.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating JS Deliver currency provider")
	}

	// SET PROVIDERS CHAIN
	govUaClient.SetNext(jsDelivrClient)
	privateClient.SetNext(govUaClient)

	return privateClient
}

func getMailClient() *mail.Client {
	mailClient, err := mail.NewClient(prodcnf.LoadProducerConfig())
	if err != nil {
		log.Fatalf("generating mail client: %v", err)
	}

	return mailClient
}

func startServer(r *gin.Engine) *http.Server {
	cnf := userconf.LoadServerConfigConfig()
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

func gracefulShutdown(ctx context.Context, scheduler *cron.Cron, server *http.Server, mailClient *mail.Client) {
	log.Println("Stopping server")

	cnf := userconf.LoadServerConfigConfig()
	waitSeconds := cnf.GracefulShutdownWaitTimeSeconds

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Duration(waitSeconds)*time.Second)
	defer shutdownCancel()

	// STOP JOBS
	scheduler.Stop()

	// STOP WEB SERVER
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}

	// STOP MAIL CLIENT
	if err := mailClient.Close(); err != nil {
		log.Printf("stopping mail client: %v", err)
	}

	select {
	case <-shutdownCtx.Done():
		log.Printf("Timeout of %d seconds\n", waitSeconds)
	}

	// STOP JOBS (hard)

	log.Println("Server exiting")
}
