package request

import (
	"errors"
	"net"
)

type AddBlackListReq struct {
	ShortUrl string `json:"shortUrl"`
	Ip       string `json:"ip"`
	UserId   uint64 `json:"userId"`
}

func (req *AddBlackListReq) Validate() error {
	if req.ShortUrl == "" {
		return errors.New("短链为空")
	}
	if req.Ip == "" {
		return errors.New("IP为空")
	}
	if net.ParseIP(req.Ip) == nil {
		return errors.New("IP无效")
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
	Ip       string `json:"ip"`
}

func (req *ListBlackListReq) Validate() error {
	if req.ShortUrl == "" {
		return errors.New("短链为空")
	}
	return nil
}
