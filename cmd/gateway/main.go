package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"

	"gateway/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs2 *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs2,
		),
	)
}

func main() {
	fmt.Println("come into main")
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	//spew.Dump(c)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	//// 进入自定义部分
	//fmt.Println("come into initClient")
	//initClient(bc.Data)
	//
	//var endpoints = []string{bc.Micro.GetLb().GetAddr()}
	//spew.Dump(endpoints)
	//
	//fmt.Println("come into discovery")
	////discovery.InitServiceDiscovery(endpoints, []string{"user-svr", "sec_kill-svr"})
	//discovery.InitServiceDiscovery(endpoints, bc.Micro.GetLb().GetDisSvrList())
	//
	//fmt.Println("come into limiter")
	//err := limiter.InitLimiter("../../configs/router.json", rdb,
	//err := limiter.InitLimiter("../../../configs/router.json", rdb,
	////	3, 10, 100)
	//if err != nil {
	//	fmt.Println("panic : ", err)
	//	panic(err)
	//}

	fmt.Println("come into wireApp")
	//app, cleanup, err := wireApp(bc.Server, bc, logger)
	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Micro, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

var (
	rdb *redis.Client
)

// 初始化连接
func initClient(cfData *conf.Data) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfData.GetRedis().GetAddr(),
		Password: cfData.GetRedis().GetPassWord(), // no password set
		DB:       0,                               // use default DB
		PoolSize: 100,                             // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}
