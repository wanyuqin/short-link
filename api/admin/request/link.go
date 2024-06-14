package request

type AddLinkReq struct {
	OriginalUrl string `json:"originalUrl"`
	ExpiredAt   int64  `json:"expiredAt"`
	UserId      uint64 `json:"userId"`
}
