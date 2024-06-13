package request

type AddLinkReq struct {
	OriginalUrl string `json:"originalUrl"`
	ExpiredAt   int64  `json:"expiredAt"`
}
