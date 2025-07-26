package auth

import (
	"context"
	"gateway/internal/conf"
	"github.com/BitofferHub/pkg/constant"
	"github.com/BitofferHub/pkg/middlewares/discovery"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type TargetHost struct {
	Host string
}

// 实现反向代理
func Forward(c *gin.Context) {
	// 请求路径为 /bitstorm/lottery-svr, 则 action = "lottery-svr"
	action, _ := c.Params.Get("action")
	// conf.Routes 是一个映射表，保存了路径名到目标后端服务的映射
	// 和定义好的 router.json 比较
	if _, ok := conf.Routes[action]; !ok {
		c.JSON(http.StatusNotFound, "") // 如果找不到对应路径, 直接返回 404
		return
	}
	route := conf.Routes[action]
	userID, _ := c.Get("userID")
	hostReverseProxy(c, c.Writer, c.Request, &route, userID.(string))
}

func hostReverseProxy(ctx context.Context, w http.ResponseWriter, req *http.Request, route *conf.Route, userID string) {
	log.Infof("redirect route: %+v\n", route)
	// 反向代理的核心, 它会在转发前修改请求信息（比如地址、请求头等），然后由代理自动转发
	// 感觉就是把 router.json 里面定义好的东西拼在一起
	director := func(req *http.Request) {
		req.Header.Set(constant.UserID, userID)
		req.Header.Set(constant.TraceID, "121321313")
		req.Header.Set("Pika-AAA", "abcdefghhhhhhhh")

		req.URL.Scheme = route.Scheme                    // 设置 http 或 https 协议
		dis := discovery.GetServiceDiscovery(route.Host) // 获取目标后端服务地址
		endpoint := dis.GetHttpEndPoint()
		u, err := url.Parse(endpoint)
		if err != nil {
			panic(err)
		}
		req.URL.Host = u.Host
		req.URL.Path = route.Uri
		log.Infof("redirect url: %+v\n", req.URL)
	}

	log.Infof("redirect host: %+v, url: %+v\n", route.Host, req)
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, req) // 由代理完成真实的转发请求
}
