package service

// Session管理服务
var Session = sessionService{}

type sessionService struct{}

const (
	// 用户信息存放在Session中的Key
	sessionKeyUser = "SessionKeyUser"
)
