package services

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"short-link/api/admin/request"
	"short-link/api/admin/resopnse"
	"short-link/ctxkit"
	"short-link/internal/link/repository"
	"short-link/internal/link/repository/db"
	"short-link/logs"
	"short-link/utils"
	"short-link/utils/timex"
	"sort"
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
	var (
		expiredAt = int64(-1)
		err       error
		logFmt    = "[LinkService][AddLink]"
	)
	// 查询是否已经存在长连接的短链
	original, err := svc.linkRepo.GetOriginal(ctx, req.OriginUrl)
	if err != nil {
		logs.Error(err, logFmt+"get original failed", zap.Any("originUrl", req.OriginUrl))
		return err
	}
	if req.ExpiredAt != "" {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", req.ExpiredAt, time.Local)
		expiredAt = t.UnixMilli()
	}
	if original != nil {
		err = errors.New("原始链接已经存在")
		logs.Error(err, logFmt, zap.Any("originUrl", req.OriginUrl))
		return err
	}
	shortLink := utils.GenerateShortLink(req.OriginUrl)
	if req.ExpiredAt != "" {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", req.ExpiredAt, time.Local)
		expiredAt = t.UnixMilli()
		if expiredAt < time.Now().UnixMilli() {
			return errors.New("过期时间不能早于当前时间")
		}
	}
	// 保存
	slLink := &db.SlLink{
		ShortUrl:  shortLink,
		OriginUrl: req.OriginUrl,
		ExpiredAt: expiredAt,
		UserId:    req.UserId,
	}
	err = svc.linkRepo.AddLink(ctx, slLink)
	logs.Error(err, logFmt+"add link failed", zap.Any("slLink", slLink))
	return err
}

func (svc *LinkService) Request(ctx context.Context, shortLink string) (string, error) {
	//TODO redis
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

	return short.OriginUrl, nil
}

func (svc *LinkService) LinkList(ctx context.Context, req *request.LinkListReq) (*resopnse.LisLinkResp, error) {
	resp := &resopnse.LisLinkResp{}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	userId := ctxkit.GetUserId(ctx)
	list, lastId, err := repository.NewLinkRepository().GetByUserId(ctx, userId, req.OriginUrl, req.LastId, req.PageSize)
	if err != nil {
		logs.Error(err, "")
		return resp, err
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt > list[j].CreatedAt
	})
	data := make([]resopnse.Link, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		slLink := resopnse.Link{
			Id:        item.ID,
			UserId:    item.UserId,
			OriginUrl: item.OriginUrl,
			ShortUrl:  item.ShortUrl,
			CreatedAt: timex.FormatDateTime(item.CreatedAt),
			UpdatedAt: timex.FormatDateTime(item.UpdatedAt),
		}
		if item.ExpiredAt > 0 {
			slLink.ExpiredAt = timex.FormatDateTime(item.ExpiredAt)
		}
		data = append(data, slLink)
	}

	resp.Data = data
	resp.LastId = lastId
	return resp, nil
}

func (svc *LinkService) DeleteLink(ctx context.Context, r *request.DelLinkReq) error {
	logFmt := "[LinkService][DeleteLink]"
	short, err := svc.linkRepo.GetByShort(ctx, r.ShortUrl)
	if err != nil {
		logs.Error(err, logFmt+"get by short link failed", zap.Any("shortLink", r.ShortUrl))
		return err
	}
	if short == nil {
		return errors.New("短链不存在")
	}
	if short.UserId != r.UserId {
		return errors.New("操作非法")
	}
	if err = svc.linkRepo.UpdateByShort(ctx, short.ShortUrl, map[string]interface{}{
		"is_del":     1,
		"updated_at": time.Now().UnixMilli(),
	}); err != nil {
		logs.Error(err, logFmt+"update by short failed")
	}

	return err
}
