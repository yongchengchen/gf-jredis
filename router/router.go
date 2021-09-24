package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/yongchengchen/gf-jredis/app/api"
	"github.com/yongchengchen/gf-jredis/app/service"
)

func init() {
	s := g.Server()

	s.Group("/jredis", func(group *ghttp.RouterGroup) {
		group.Middleware(
			service.Middleware.Ctx,
			service.Middleware.CORS,
			service.Middleware.InnerAuth,
		)

		group.POST("/set", api.Jredis.SetValue)
		group.GET("/get/:key", api.Jredis.GetValue)

		group.POST("/hset", api.Jredis.HSetValue)
		group.GET("/hget/:key/:field", api.Jredis.HGetValue)

		group.POST("/push/:direction", api.Jredis.Push)
		group.POST("/maxlpush/:limit", api.Jredis.LimitLPush)

		group.POST("/push/:direction", api.Jredis.Push)
		group.POST("/maxlpush/:limit", api.Jredis.LimitLPush)

		group.GET("/lrange/:key/:from/:to", api.Jredis.LRange)
		group.GET("/llen/:key", api.Jredis.LLen)
		group.GET("/hkeys/:key", api.Jredis.HKeys)

		group.POST("/expire/:key/:ttl", api.Jredis.Expire)
		group.GET("/ttl/:key", api.Jredis.Ttl)
		group.DELETE("/del/:key", api.Jredis.Del)
		group.GET("/info", api.Jredis.Info)
	})

	// 分组路由注册方式
	// s.Group("/", func(group *ghttp.RouterGroup) {
	// 	group.Middleware(
	// 		service.Middleware.Ctx,
	// 		service.Middleware.CORS,
	// 	)
	// 	// group.ALL("/chat", api.Chat)
	// 	group.ALL("/user", api.User)
	// 	group.Group("/", func(group *ghttp.RouterGroup) {
	// 		group.Middleware(service.Middleware.Auth)
	// 		group.ALL("/user/profile", api.User.Profile)
	// 	})
	// })
}
