package resopnse

type LisLinkResp struct {
	Data   []Link `json:"data"`
	LastId uint64 `json:"lastId"`
}

type Link struct {
	Id        uint64 `json:"id"`
	UserId    uint64 `json:"userId"`
	OriginUrl string `json:"originUrl"`
	ShortUrl  string `json:"shortUrl"`
	ExpiredAt string `json:"expiredAt"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
