package services

import (
	"context"
	"errors"
	"short-link/api/request"
	"short-link/internal/link/repository"
	"short-link/internal/link/repository/db"
	"short-link/utils"
)

type LinkService struct {
	linkRepo repository.ILinkRepository
}

func NewLinkService() *LinkService {
	return &LinkService{
		linkRepo: repository.NewLinkRepository(),
	}
}

func (svc *LinkService) AddLink(ctx context.Context, req *request.AddLinkReq) error {
	if req.OriginalUrl == "" {
		return errors.New("原链接为空")
	}
	// 查询是否已经存在长连接的短链
	original, err := svc.linkRepo.GetOriginal(ctx, req.OriginalUrl)
	if err != nil {
		return err
	}
	if original != nil {
		return errors.New("原始链接已经存在")
	}
	shortLink := utils.GenerateShortLink(req.OriginalUrl)

	// 碰撞校验

	// 保存
	err = svc.linkRepo.AddLink(ctx, &db.SlLink{
		ShortUrl:  shortLink,
		OriginUrl: req.OriginalUrl,
		ExpiredAt: req.ExpiredAt,
		UserId:    req.UserId,
	})
	return err
}
