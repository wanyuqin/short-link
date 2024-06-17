package services

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"short-link/api/admin/request"
	"short-link/ctxkit"
	"short-link/internal/link/repository"
	"short-link/internal/link/repository/db"
	"short-link/logs"
	"short-link/utils"
	"time"
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
	var expiredAt = int64(-1)
	userId := ctxkit.GetUserId(ctx)
	if userId == 0 {
		return errors.New("未登录")
	}
	req.UserId = userId
	if req.OriginUrl == "" {
		return errors.New("原链接为空")
	}
	// 查询是否已经存在长连接的短链
	original, err := svc.linkRepo.GetOriginal(ctx, req.OriginUrl)
	if err != nil {
		return err
	}
	if original != nil {
		return errors.New("原始链接已经存在")
	}
	shortLink := utils.GenerateShortLink(req.OriginUrl)
	if req.ExpiredAt != "" {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", req.ExpiredAt, time.Local)
		expiredAt = t.Unix()
	}
	// 保存
	err = svc.linkRepo.AddLink(ctx, &db.SlLink{
		ShortUrl:  shortLink,
		OriginUrl: req.OriginUrl,
		ExpiredAt: expiredAt,
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
