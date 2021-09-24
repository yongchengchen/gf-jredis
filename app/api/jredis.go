package api

import (
	"strconv"
	"strings"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/yongchengchen/gf-jredis/app/model"
	"github.com/yongchengchen/gf-jredislibrary/response"
)

// 用户API管理对象
var Jredis = new(jredisApi)

type jredisApi struct{}

// redis Set command
func (a *jredisApi) SetValue(r *ghttp.Request) {
	var (
		data *model.JredisApiKeyValuePairReq
	)
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 400, err.Error())
	}

	g.Redis().Do("SET", data.Key, data.Data)

	response.JsonExit(r, 201, data.Key, data.Data)
}

// redis Get command
func (a *jredisApi) GetValue(r *ghttp.Request) {
	var (
		err    error
		result *gvar.Var
	)

	var key = r.GetString("key")
	result, err = g.Redis().DoVar("GET", key)
	if err != nil {
		response.JsonExit(r, 500, err.Error())
	}
	response.JsonExit(r, 200, "success", result.String())
}

// redis HSet command
func (a *jredisApi) HSetValue(r *ghttp.Request) {
	var (
		data *model.JredisApiHSetValueReq
	)
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 400, err.Error())
	}

	g.Redis().Do("HSET", data.Key, data.Field, data.Data)

	response.JsonExit(r, 201, data.Key, data.Data)
}

// redis HGet and HGETALL command
func (a *jredisApi) HGetValue(r *ghttp.Request) {
	var (
		err    error
		result *gvar.Var
	)

	var key = r.GetString("key")
	var field = r.GetString("field")

	if field == "all" {
		result, err = g.Redis().DoVar("HGETALL", key)
	} else {
		result, err = g.Redis().DoVar("HGET", key, field)
	}

	if err != nil {
		response.JsonExit(r, 500, err.Error())
	}
	response.JsonExit(r, 200, "success", result.String())
}

// redis LPUSH/RPUSH command
func (a *jredisApi) Push(r *ghttp.Request) {
	var (
		data *model.JredisApiKeyValuePairReq
	)
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 400, err.Error())
	}

	cmd := ""
	var direction = r.GetString("direction") //l or r push
	switch strings.ToUpper(direction) {
	case "L":
		cmd = "LPUSH"
	case "LPUSH":
		cmd = "LPUSH"
	case "R":
		cmd = "RPUSH"
	case "RPUSH":
		cmd = "RPUSH"
	}
	if len(cmd) > 0 {
		_, err := g.Redis().Do(cmd, data.Key, data.Data)
		if err != nil {
			response.JsonExit(r, 500, "Fail to push")
		}
	} else {
		response.JsonExit(r, 400, "wrong push direction", strings.ToUpper(direction))
		return
	}

	response.JsonExit(r, 201, data.Key, data.Data)
}

// redis LPUSH command
func (a *jredisApi) LimitLPush(r *ghttp.Request) {
	var limit = r.GetString("limit") //l or r push
	max, err := strconv.Atoi(limit)
	if err != nil {
		response.JsonExit(r, 400, "wrong max list length")
	}

	var (
		data *model.JredisApiKeyValuePairReq
	)
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 400, err.Error())
	}
	_, err = g.Redis().DoVar("LPUSH", data.Key, data.Data)
	if err != nil {
		response.JsonExit(r, 500, "Fail to push")
	}

	if max > 1 {
		g.Redis().DoVar("LTRIM", data.Key, 0, max-1)
	}

	response.JsonExit(r, 201, data.Key, data.Data)
}

func (a *jredisApi) LLen(r *ghttp.Request) {
	var (
		err    error
		result *gvar.Var
	)

	var key = r.GetString("key")
	result, err = g.Redis().DoVar("LLEN", key)
	if err != nil {
		response.JsonExit(r, 500, err.Error())
	}
	response.JsonExit(r, 200, "success", result.Int32())
}

func (a *jredisApi) HKeys(r *ghttp.Request) {
	var (
		err    error
		result *gvar.Var
	)

	var key = r.GetString("key")
	if key == "gjredis_auth" {
		response.JsonExit(r, 200, "success")
		return
	}

	result, err = g.Redis().DoVar("HKeys", key)
	if err != nil {
		response.JsonExit(r, 500, err.Error())
	}
	response.JsonExit(r, 200, "success", result.Interfaces())
}

func (a *jredisApi) LRange(r *ghttp.Request) {
	var (
		err    error
		result *gvar.Var
	)

	var key = r.GetString("key")
	var from = r.GetString("from")
	var to = r.GetString("to")

	result, err = g.Redis().DoVar("LRANGE", key, from, to)

	if err != nil {
		response.JsonExit(r, 500, err.Error())
	}
	response.JsonExit(r, 200, "success", result.Interfaces())
}

// redis Del command
func (a *jredisApi) Del(r *ghttp.Request) {
	var key = r.GetString("key")
	if key == "gjredis_auth" {
		response.JsonExit(r, 500, "Can't delete.")
		return
	}

	simpleKeyCommand("DEL", key, nil, r)
}

// redis Expire command
func (a *jredisApi) Expire(r *ghttp.Request) {
	var key = r.GetString("key")
	var ttl = r.GetString("ttl")
	if key == "gjredis_auth" {
		response.JsonExit(r, 500, "Reserve Key")
		return
	}

	simpleKeyCommand("EXPIRE", key, ttl, r)
}

// redis Ttl command
func (a *jredisApi) Ttl(r *ghttp.Request) {
	var key = r.GetString("key")
	if key == "gjredis_auth" {
		response.JsonExit(r, 500, "Reserve Key")
		return
	}

	simpleKeyCommand("TTL", key, nil, r)
}

func (a *jredisApi) Info(r *ghttp.Request) {
	var (
		data *model.JredisApiInfoFieldReq
	)
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 400, err.Error())
	}

	var (
		err    error
		result *gvar.Var
	)

	if len(data.Field) > 0 {
		result, err = g.Redis().DoVar("info", data.Field)
	} else {
		result, err = g.Redis().DoVar("info")
	}
	if err != nil {
		response.JsonExit(r, 500, err.Error())
	}
	response.JsonExit(r, 200, "success", result.Interfaces())
}

func simpleKeyCommand(command string, key string, value interface{}, r *ghttp.Request) {
	var (
		err    error
		result *gvar.Var
	)
	if value == nil {
		result, err = g.Redis().DoVar(command, key)
	} else {
		result, err = g.Redis().DoVar(command, key, value)
	}
	if err != nil {
		response.JsonExit(r, 500, err.Error())
	}

	response.JsonExit(r, 200, "success", result.Interfaces())
}
