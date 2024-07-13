package request

import (
	"errors"
	"net"
	"short-link/internal/consts"
)

type AddBlackListReq struct {
	ShortURL string `json:"shortUrl"`
	IP       string `json:"ip"`
	UserID   uint64 `json:"userId"`
}

func (req *AddBlackListReq) Validate() error {
	if req.ShortURL == "" {
		return consts.ErrShortURLIsEmpty
	}
	if req.IP == "" {
		return consts.ErrIPIsEmpty
	}
	if net.ParseIP(req.IP) == nil {
		return consts.ErrIpIsInvalid
	}
	return nil
}

type DeleteBlackListReq struct {
	Id uint64 `json:"id"`
}

func (req *DeleteBlackListReq) Validate() error {
	if req.Id == 0 {
		return errors.New("id 为空")
	}
	return nil
}

type ListBlackListReq struct {
	PageReq
	ShortUrl string `json:"shortUrl"`
	IP       string `json:"ip"`
	Status   int    `json:"status"`
}

func (req *ListBlackListReq) Validate() error {
	if req.ShortUrl == "" {
		return consts.ErrShortURLIsEmpty
	}
	return nil
}

type UpdateBlackListReq struct {
	ID       uint64 `json:"id"`
	Ip       string `json:"ip"`
	Status   int    `json:"status"`
	ShortURL string `json:"shortUrl"`
}

func (req *UpdateBlackListReq) Validate() error {
	if req.ID == 0 {
		return errors.New("id 为空")
	}
	if req.ShortURL == "" {
		return consts.ErrShortURLIsEmpty
	}
	return nil
}
