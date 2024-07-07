package request

import (
	"errors"
	"short-link/internal/consts"
	"time"
)

type PageReq struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type AddLinkReq struct {
	OriginURL string `json:"originUrl"`
	ExpiredAt string `json:"expiredAt"`
	UserID    uint64 `json:"userId"`
}

func (req *AddLinkReq) Validate() error {
	if req.OriginURL == "" {
		return consts.ErrURLIsEmpty
	}
	if req.ExpiredAt != "" {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", req.ExpiredAt, time.Local)
		if t.UnixMilli() < time.Now().UnixMilli() {
			return errors.New("过期时间不能早于当前时间")
		}
	}
	return nil
}

type LinkListReq struct {
	PageReq
	OriginUrl string `json:"originUrl"`
}

type DelLinkReq struct {
	ShortUrl string `json:"shortUrl"`
	UserId   uint64 `json:"userId"`
}

func (req *DelLinkReq) Validate() error {
	if req.ShortUrl == "" {
		return consts.ErrShortURLIsEmpty
	}
	return nil
}
