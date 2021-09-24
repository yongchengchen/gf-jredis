package model

type JredisApiKeyValuePairReq struct {
	Key  string      `json:"key" v:"required#Redis Set Key can not be empty"`
	Data interface{} `json:"data" v:""`
}

type JredisApiHSetValueReq struct {
	Key   string      `json:"key" v:"required#Redis HSet Key can not be empty"`
	Field string      `json:"field" v:"required#Redis HSet Field can not be empty"`
	Data  interface{} `json:"data" v:""`
}

type JredisApiOnlyKeyFieldReq struct {
	Key   string `json:"key" v:"required#Redis Key can not be empty"`
	Field string `json:"field" v:""`
}

type JredisApiInfoFieldReq struct {
	Field string `json:"field" v:""`
}
