package consts

import "errors"

const TokenSecret = "172afec60b54adfd982b710ed4b54b83"

const (
	HTTPScheme  = "http"
	HTTPSScheme = "https"
)

var (
	ErrURLIsEmpty  = errors.New("地址为空")
	ErrHostIsEmpty = errors.New("host为空")

	ErrShortURLExpired = errors.New("短链过期")
	ErrShortURLIsEmpty = errors.New("短链为空")

	ErrIPBlocked   = errors.New("IP被拦截")
	ErrIPIsEmpty   = errors.New("IP为空")
	ErrIpIsInvalid = errors.New("IP无效")

	ErrSchemeIsEmpty = errors.New("协议为空")
	ErrSchemeInvalid = errors.New("协议非法")
)

const (
	IPStatusDisable = 0 // 禁用
	IPStatusActive  = 1 // 启动
)
