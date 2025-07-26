package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
//var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)

// ProviderSet 使用新封装好的中间件实现依赖注入
var ProviderSet = wire.NewSet(NewHTTPServer2, NewGRPCServer)
