package resopnse

type LisLinkResp struct {
	Data  []Link `json:"data"`
	Total int64  `json:"total"`
}

type Link struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"userId"`
	OriginURL string `json:"originUrl"`
	ShortURL  string `json:"shortUrl"`
	ExpiredAt string `json:"expiredAt"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
