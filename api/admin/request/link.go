package request

type AddLinkReq struct {
	OriginUrl string `json:"originUrl"`
	ExpiredAt string `json:"expiredAt"`
	UserId    uint64 `json:"userId"`
}

type LinkListReq struct {
	LastId    uint64 `json:"lastId"`
	OriginUrl string `json:"originUrl"`
	PageSize  uint64 `json:"pageSize"`
}
