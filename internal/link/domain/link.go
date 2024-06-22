package domain

type ShorUrl struct {
	OriginUrl string `json:"originUrl,omitempty"`
	ShorUrl   string `json:"shorUrl,omitempty"`
	ExpiredAt int64  `json:"expiredAt,omitempty"`
}
