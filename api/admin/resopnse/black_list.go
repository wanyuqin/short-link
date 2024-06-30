package resopnse

type ListBlackListResp struct {
	Data  []BlackList `json:"data,omitempty"`
	Total int64       `json:"total,omitempty"`
}

type BlackList struct {
	Id        uint64 `json:"id,omitempty"`
	ShortUrl  string `json:"shortUrl,omitempty"`
	Ip        string `json:"ip,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}
