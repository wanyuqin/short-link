package services

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/url"
	"short-link/api/admin/request"
	"short-link/api/admin/resopnse"
	"short-link/ctxkit"
	"short-link/database/cache"
	"short-link/internal/consts"
	"short-link/internal/link/domain"
	"short-link/internal/link/repository"
	"short-link/internal/link/repository/db"
	"short-link/logs"
	"short-link/utils"
	"short-link/utils/timex"
	"sort"
	"strconv"
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
	if err = svc.validateOriginUrl(req.OriginUrl); err != nil {
		logs.Error(err, logFmt+"validate origin url failed")
		return err
	}
	// 查询是否已经存在长连接的短链
	original, err := svc.linkRepo.GetOriginal(ctx, req.OriginUrl, req.UserId)
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
	// TODO 完善生成逻辑，减少hash冲突
	shortLink := utils.GenerateShortLink(req.OriginUrl + strconv.FormatUint(req.UserId, 10))
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
	if err != nil {
		logs.Error(err, logFmt+"add link failed", zap.Any("slLink", slLink))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("原链接重复")
		}
	}

	return err
}

func (svc *LinkService) Request(ctx context.Context, shortLink string) (string, error) {
	var (
		redisKey  = fmt.Sprintf(consts.RedisKeyShorUrl, shortLink)
		originUrl string
		shortUrl  domain.ShorUrl
	)

	rdb := cache.NewRedisTool(ctx)
	err := rdb.AutoFetch(ctx, redisKey, 0, &shortUrl, func(ctx context.Context) (interface{}, error) {
		short, err := svc.linkRepo.GetByShort(ctx, shortLink)
		if err != nil {
			return "", err
		}

		if short == nil {
			return "", errors.New("record not found")
		}
		ds := domain.ShorUrl{
			OriginUrl: short.OriginUrl,
			ShorUrl:   short.ShortUrl,
			ExpiredAt: short.ExpiredAt,
		}
		return ds, nil
	})
	if err != nil {
		logs.Error(err, "[LinkService] auto fetch failed")
		return "", err
	}

	if shortUrl.ExpiredAt > 0 && shortUrl.ExpiredAt < time.Now().UnixMilli() {
		return "", errors.New("record not found")
	}
	originUrl = shortUrl.OriginUrl

	return originUrl, nil
}

func (svc *LinkService) LinkList(ctx context.Context, req *request.LinkListReq) (*resopnse.LisLinkResp, error) {
	resp := &resopnse.LisLinkResp{}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	userId := ctxkit.GetUserId(ctx)
	list, lastId, err := repository.NewLinkRepository().PageUserLink(ctx, userId, req.OriginUrl, req.LastId, req.PageSize)
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
	if err = svc.linkRepo.DeleteByShort(ctx, short.ShortUrl); err != nil {
		logs.Error(err, logFmt+"delete by short failed")
		return err
	}

	return err
}

func (svc *LinkService) validateOriginUrl(originUrl string) error {
	if originUrl == "" {
		return consts.ErrUrlIsEmpty
	}
	parseUrl, err := url.Parse(originUrl)
	if err != nil {
		return err
	}
	if parseUrl.Scheme == "" {
		return consts.ErrSchemeIsEmpty
	}
	if parseUrl.Scheme != consts.HttpScheme && parseUrl.Scheme != consts.HttpsScheme {
		return consts.ErrSchemeInvalid
	}
	if parseUrl.Host == "" {
		return consts.ErrHostIsEmpty
	}
	return nil
}
