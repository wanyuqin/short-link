package domain

import (
	"time"
)

type ShorUrl struct {
	OriginUrl string `json:"originUrl,omitempty"`
	ShorUrl   string `json:"shorUrl,omitempty"`
	ExpiredAt int64  `json:"expiredAt,omitempty"`
}

func (s *ShorUrl) IsValid() bool {
	if s.ExpiredAt > 0 && s.ExpiredAt < time.Now().UnixMilli() {
		return false
	}
	return true
}
