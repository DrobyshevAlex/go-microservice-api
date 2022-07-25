package src

import (
	"context"
	"log"
	"main/src/config"
	"main/src/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	go_amqp_lib "github.com/lan143/go-amqp-lib"
	go_healthcheck_lib "github.com/lan143/go-healthcheck-lib"
)

type Application struct {
	amqpClient  *go_amqp_lib.Client
	config      *config.Config
	healthCheck *go_healthcheck_lib.HealthCheck
	webServer   *http.Server

	wg             sync.WaitGroup
	sigs           chan os.Signal
	wsShutdownChan chan struct{}
}

func (a *Application) Init(ctx context.Context) error {
	log.Printf("application: init")

	a.sigs = make(chan os.Signal, 1)
	a.wsShutdownChan = make(chan struct{}, 1)

	err := a.config.Init()
	if err != nil {
		return err
	}

	signal.Notify(a.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	a.amqpClient.Init(a.config.GetAmqpConfig(), 10, &a.wg)

	a.webServer.Init(
		a.config.GetWebServerAddress(),
		&a.wg,
		a.config.IsDebug(),
		a.wsShutdownChan,
	)
	a.healthCheck.Init(a.config.GetHealthCheckAddress(), &a.wg)
	a.healthCheck.AddReadinessProbe(a.amqpClient)
	a.healthCheck.AddReadinessProbe(a.webServer)

	return nil
}

func (a *Application) Run(ctx context.Context) error {
	log.Printf("application.run: start")

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go a.processSignals(cancelFunc)

	a.healthCheck.Run(cancelCtx)

	err := a.amqpClient.Run(cancelCtx)
	if err != nil {
		return err
	}

	err = a.webServer.Run(cancelCtx)
	if err != nil {
		return err
	}

	go func() {
		<-a.wsShutdownChan
		a.amqpClient.Shutdown(ctx)
	}()

	log.Println("application.run: running")

	a.wg.Wait()

	log.Println("application: graceful shutdown.")

	return nil
}

func (a *Application) processSignals(cancelFun context.CancelFunc) {
	select {
	case <-a.sigs:
		log.Println("application: received shutdown signal from OS")
		cancelFun()
		break
	}
}

func NewApplication(
	amqpClient *go_amqp_lib.Client,
	config *config.Config,
	healthCheck *go_healthcheck_lib.HealthCheck,
	webServer *http.Server,
) *Application {
	return &Application{
		amqpClient:  amqpClient,
		config:      config,
		healthCheck: healthCheck,
		webServer:   webServer,
	}
}
