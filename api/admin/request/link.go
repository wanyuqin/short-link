package request

type AddLinkReq struct {
	OriginUrl string `json:"originUrl"`
	ExpiredAt string `json:"expiredAt"`
	UserId    uint64 `json:"userId"`
}

type LinkListReq struct {
	LastId    uint64 `json:"lastId"`
	PageSize  uint64 `json:"pageSize"`
	OriginUrl string `json:"originUrl"`
}
