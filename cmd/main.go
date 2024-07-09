package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"genesis-currency-api/internal/db"
	dbconf "genesis-currency-api/internal/db/config"
	"genesis-currency-api/internal/job"
	"genesis-currency-api/internal/middleware"
	currencymodule "genesis-currency-api/internal/module/currency"
	ratecdnjsdelivr "genesis-currency-api/internal/module/currency/api/external/rater/cdnjsdelivr"
	rategovua "genesis-currency-api/internal/module/currency/api/external/rater/gov_ua"
	rateprivate "genesis-currency-api/internal/module/currency/api/external/rater/private"
	currencyhand "genesis-currency-api/internal/module/currency/api/handler"
	currencyconf "genesis-currency-api/internal/module/currency/config"
	currencyserv "genesis-currency-api/internal/module/currency/service"
	emailmodule "genesis-currency-api/internal/module/email"
	emailhand "genesis-currency-api/internal/module/email/api/handler"
	emailconf "genesis-currency-api/internal/module/email/config"
	usermodule "genesis-currency-api/internal/module/user"
	userhand "genesis-currency-api/internal/module/user/api/handler"
	userconf "genesis-currency-api/internal/server/config"
	"genesis-currency-api/pkg/common/envs"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	envs.Init()

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
	emailModule := emailmodule.Init(userModule.Service, currencyModule.Service, emailconf.LoadEmailServiceConfig())

	// ENGINE
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	// JOBS
	scheduler := job.StartAllJobs(
		job.GetUpdateCurrencyJob(ctx, currencyModule.Service),
		job.GetSendEmailsJob(ctx, emailModule.Service),
	)

	// HANDLERS
	currencyhand.RegisterRoutes(r, currencyModule.Handler)
	userhand.RegisterRoutes(r, userModule.Handler)
	emailhand.RegisterRoutes(r, emailModule.Handler)

	// START SERVER
	server := startServer(r)
	waitServerWorking()

	// STOP SERVER
	gracefulShutdown(ctx, scheduler, server)
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

func gracefulShutdown(ctx context.Context, scheduler *cron.Cron, server *http.Server) {
	log.Println("Stopping server")

	cnf := userconf.LoadServerConfigConfig()
	waitSeconds := cnf.GracefulShutdownWaitTimeSeconds

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Duration(waitSeconds)*time.Second)
	defer shutdownCancel()

	// STOP JOBS
	scheduler.Stop()

	// STOP WEB SERVER
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server shutdown:", err)
	}

	select {
	case <-shutdownCtx.Done():
		log.Printf("Timeout of %d seconds\n", waitSeconds)
	}

	// STOP JOBS (hard)

	log.Println("Server exiting")
}
