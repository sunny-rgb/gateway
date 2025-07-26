//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/BitofferHub/gateway/internal/biz"
	"github.com/BitofferHub/gateway/internal/conf"
	"github.com/BitofferHub/gateway/internal/data"
	"github.com/BitofferHub/gateway/internal/server"
	"github.com/BitofferHub/gateway/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Micro, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
