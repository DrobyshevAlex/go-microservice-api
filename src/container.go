package src

import (
	go_amqp_lib "github.com/lan143/go-amqp-lib"
	go_healthcheck_lib "github.com/lan143/go-healthcheck-lib"

	"main/src/config"
	"main/src/http"
	v1_0_0 "main/src/http/controllers/v1.0.0"
	"main/src/services"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	processError(container.Provide(NewApplication))
	processError(container.Provide(http.NewServer))
	processError(container.Provide(config.NewConfig))

	// 3rd party
	processError(container.Provide(go_amqp_lib.NewClient))
	processError(container.Provide(go_healthcheck_lib.NewHealthCheck))

	// Controllers
	// v1.0.0
	processError(container.Provide(v1_0_0.NewApiController))
	processError(container.Provide(v1_0_0.NewUserController))

	// Services
	processError(container.Provide(services.NewUserService))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
