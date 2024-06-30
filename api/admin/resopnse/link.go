package resopnse

type LisLinkResp struct {
	Data  []Link `json:"data"`
	Total int64  `json:"total"`
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
