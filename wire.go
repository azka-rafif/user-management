//go:build wireinject
// +build wireinject

package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/domain/auth"
	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/evermos/boilerplate-go/transport/http"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/router"
	"github.com/google/wire"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistences = wire.NewSet(
	infras.ProvideMySQLConn,
)

var domainAuth = wire.NewSet(
	auth.ProvideAuthServiceImpl,
	wire.Bind(new(auth.AuthService), new(*auth.AuthServiceImpl)),
	auth.ProvideAuthRepositoryMySQL,
	wire.Bind(new(auth.AuthRepository), new(*auth.AuthRepositoryMySQL)),
)

var domainProduct = wire.NewSet(
	product.ProvideProductServiceImpl,
	wire.Bind(new(product.ProductService), new(*product.ProductServiceImpl)),
	product.ProvideProductRepositoryMySQL,
	wire.Bind(new(product.ProductRepository), new(*product.ProductRepositoryMySQL)),
)

var domainCart = wire.NewSet(
	cart.ProvideCartServiceImpl,
	wire.Bind(new(cart.CartService), new(*cart.CartServiceImpl)),
	cart.ProvideCartRepositoryMySQL,
	wire.Bind(new(cart.CartRepository), new(*cart.CartRepositoryMySQL)),
)

var domainOrder = wire.NewSet(
	order.ProvideOrderServiceImpl,
	wire.Bind(new(order.OrderService), new(*order.OrderServiceImpl)),
	order.ProvideOrderRepositoryMySQL,
	wire.Bind(new(order.OrderRepository), new(*order.OrderRepositoryMySQL)),
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainAuth, domainProduct, domainCart, domainOrder,
)

var authMiddleware = wire.NewSet(
	middleware.ProvideAuthentication,
	middleware.ProvideJwtAuthentication,
)

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "AuthHandler", "ProductHandler", "CartHandler", "OrderHandler"),
	handlers.ProvideAuthHandler,
	handlers.ProvideCartHandler,
	handlers.ProvideOrderHandler,
	handlers.ProvideProductHandler,
	router.ProvideRouter,
)

// Wiring for everything.
func InitializeService() *http.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// middleware
		authMiddleware,
		// domains
		domains,
		// routing
		routing,
		// selected transport layer
		http.ProvideHTTP)
	return &http.HTTP{}
}
