package consts

import "errors"

const TokenSecret = "172afec60b54adfd982b710ed4b54b83"

const (
	HttpScheme  = "http"
	HttpsScheme = "https"
)

var (
	ErrUrlIsEmpty      = errors.New("url is empty")
	ErrSchemeIsEmpty   = errors.New("scheme is empty")
	ErrSchemeInvalid   = errors.New("scheme invalid")
	ErrHostIsEmpty     = errors.New("host is empty")
	ErrShortUrlExpired = errors.New("short url expired")
	ErrIpBlocked       = errors.New("ip blocked")
)
