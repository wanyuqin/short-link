package consts

import "errors"

const TokenSecret = "172afec60b54adfd982b710ed4b54b83"

const (
	HTTPScheme  = "http"
	HTTPSScheme = "https"
)

var (
	ErrURLIsEmpty      = errors.New("url is empty")
	ErrSchemeIsEmpty   = errors.New("scheme is empty")
	ErrSchemeInvalid   = errors.New("scheme invalid")
	ErrHostIsEmpty     = errors.New("host is empty")
	ErrShortURLExpired = errors.New("short url expired")
	ErrIPBlocked       = errors.New("IP blocked")
)
