package main

import (
	"fmt"
	"log"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/container"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/src/handlers"
	"github.com/evermos/boilerplate-go/src/repositories"
	"github.com/evermos/boilerplate-go/src/services"

	"github.com/evermos/boilerplate-go/docs" // swagger Docs
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	serviceName     = "Evermos/ExampleService"
	serviceVersion  = "0.0.1"
	environtmentDev = "development"
)

var (
	db     *infras.MysqlConn
	config *configs.Config
)

type Router struct {
	ExampleHandler *handlers.ExampleHandler `inject:"handler.example"`
}

func registry() *container.ServiceRegistry {
	c := container.NewContainer()
	config = configs.Get()
	db = &infras.MysqlConn{Write: infras.WriteMysqlDB(*config), Read: infras.ReadMysqlDB(*config)}
	c.Register("config", config)
	c.Register("db", db)

	// Repository
	c.Register("repository.example", new(repositories.ExampleRepository))

	// Service
	c.Register("service.example", new(services.ExampleService))

	// Handler
	c.Register("handler.example", new(handlers.ExampleHandler))

	return c
}

func Routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	c := registry()
	router := Router{}

	c.Register("router", &router)

	if err := c.Start(); err != nil {
		log.Fatalln(err)
	}
	if config.Env == environtmentDev {
		docs.SwaggerInfo.Title = serviceName
		docs.SwaggerInfo.Version = serviceVersion
		swaggerURL := fmt.Sprintf("%s/swagger/doc.json", config.AppURL)
		mux.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))
	}

	router.ExampleHandler.Router(mux)

	return mux
}