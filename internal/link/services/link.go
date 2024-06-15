package services

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"short-link/api/admin/request"
	"short-link/internal/link/repository"
	"short-link/internal/link/repository/db"
	"short-link/logs"
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

	// 保存
	err = svc.linkRepo.AddLink(ctx, &db.SlLink{
		ShortUrl:  shortLink,
		OriginUrl: req.OriginalUrl,
		ExpiredAt: req.ExpiredAt,
		UserId:    req.UserId,
	})
	return err
}

func (svc *LinkService) Request(ctx context.Context, shortLink string) (string, error) {

	//rdb := cache.NewRedisTool(ctx)
	//rdb.AutoFetch(ctx,)
	short, err := svc.linkRepo.GetByShort(ctx, shortLink)
	if err != nil {
		logs.Error(err, "get by short link failed", zap.Any("shortLink", shortLink))
		return "", err
	}

	if short == nil {
		logs.Error(err, "get by short link not found", zap.Any("shortLink", shortLink))
		return "", errors.New("record not found")
	}

	return short.OriginalUrl, nil
}
