package service

import (
	"net/http"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/yongchengchen/gf-jredis/app/model"
)

// 中间件管理服务
var Middleware = middlewareService{}

type middlewareService struct{}

// 自定义上下文对象
func (s *middlewareService) Ctx(r *ghttp.Request) {
	// 初始化，务必最开始执行
	customCtx := &model.Context{
		Session: r.Session,
	}
	Context.Init(r, customCtx)
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// 允许接口跨域请求
func (s *middlewareService) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

//允许内部请求
func (s *middlewareService) InnerAuth(r *ghttp.Request) {
	var authCode = r.Header.Get("Authorization")
	if len(authCode) <= 0 {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}

	var (
		err    error
		result *gvar.Var
	)

	result, err = g.Redis().DoVar("HGET", "gjredis_auth", authCode)
	if err != nil {
		r.Response.Write(err.Error())
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}

	if len(result.String()) > 0 {
		r.Middleware.Next()
	} else {
		r.Response.Write("gjredis_auth fail   \n")
		r.Response.WriteStatus(http.StatusForbidden)
	}
}
