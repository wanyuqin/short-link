package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"short-link/api/admin/request"
	"short-link/api/admin/resopnse"
	"short-link/ctxkit"
	"short-link/internal/consts"
	"short-link/internal/link/domain"
	"short-link/internal/link/repository"
	"short-link/internal/link/repository/db"
	"short-link/logs"
	"short-link/pkg/bus"
	"short-link/utils"
	"short-link/utils/gox"
	"short-link/utils/timex"
	"sort"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LinkService struct {
	linkRepo      repository.ILinkRepository
	blackListRepo repository.IBlackListRepository
}

func NewLinkService() *LinkService {
	return &LinkService{
		linkRepo:      repository.NewLinkRepository(),
		blackListRepo: repository.NewBlackListRepository(),
	}
}

func (svc *LinkService) AddLink(ctx context.Context, req *request.AddLinkReq) error {
	var (
		expiredAt = int64(-1)
		err       error
		logFmt    = "[LinkService][AddLink]"
	)

	if err = svc.validateOriginUrl(req.OriginURL); err != nil {
		logs.Error(err, logFmt+"validate origin url failed")
		return err
	}
	// 查询是否已经存在长连接的短链
	original, err := svc.linkRepo.GetOriginal(ctx, req.OriginURL, req.UserID)
	if err != nil {
		logs.Error(err, logFmt+"get original failed", zap.Any("originUrl", req.OriginURL))
		return err
	}
	if original != nil {
		err = errors.New("原始链接已经存在")
		logs.Error(err, logFmt, zap.Any("originUrl", req.OriginURL))
		return err
	}

	if req.ExpiredAt != "" {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", req.ExpiredAt, time.Local)
		if t.UnixMilli() < time.Now().UnixMilli() {
			return errors.New("过期时间不能早于当前时间")
		}

		expiredAt = t.UnixMilli()
	}

	shortLink := utils.GenerateShortLink(req.OriginURL + strconv.FormatUint(req.UserID, 10))
	// 保存
	slLink := &db.SlLink{
		ShortURL:  shortLink,
		OriginURL: req.OriginURL,
		ExpiredAt: expiredAt,
		UserID:    req.UserID,
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

func (svc *LinkService) Request(ctx context.Context, shortUrl string) (string, error) {
	var (
		originUrl  string
		blackList  domain.BlackLists
		shortUrlDo domain.ShorUrl
		logFmt     = "[LinkService][Request]"
	)
	ip := ctxkit.GetIP(ctx)
	wg := gox.NewErrorWaitGroup()

	wg.RunSafe(ctx, func(ctx context.Context) error {
		suCache, err := svc.linkRepo.GetByShortWithCache(ctx, shortUrl)
		if err != nil {
			logs.Error(err, logFmt+"get short url with cache failed", zap.Any("shortUrl", shortUrl))
			return err
		}
		if suCache != nil {
			shortUrlDo.ShorURL = suCache.ShortURL
			shortUrlDo.OriginURL = suCache.OriginURL
			shortUrlDo.ExpiredAt = suCache.ExpiredAt
		}
		return nil
	})

	wg.RunSafe(ctx, func(ctx context.Context) error {
		list, err := svc.blackListRepo.GetBlackListWithCache(ctx, shortUrl)
		blackList = list
		if err != nil {
			logs.Error(err, logFmt+"get black list  with cache failed", zap.Any("shortUrl", shortUrl))
		}
		return err
	})

	err := wg.Wait()
	if err != nil {
		return "", err
	}

	for _, item := range blackList {
		if ip == item.IP && item.Status == consts.IPStatusActive {
			logs.Warn(logFmt+"IP is blocked", zap.Any("shortUrl", shortUrl), zap.Any("IP", ip))
			return "", consts.ErrIPBlocked
		}
	}

	if !shortUrlDo.IsValid() {
		return "", consts.ErrShortURLExpired
	}

	originUrl = shortUrlDo.OriginURL

	return originUrl, nil
}

func (svc *LinkService) LinkList(ctx context.Context, req *request.LinkListReq) (*resopnse.LisLinkResp, error) {
	logFmt := "[LinkService][LinkList]"
	resp := &resopnse.LisLinkResp{}
	userId := ctxkit.GetUserId(ctx)
	list, total, err := repository.NewLinkRepository().PageUserLink(ctx, userId, req.OriginUrl, req.Page, req.PageSize)
	if err != nil {
		logs.Error(err, logFmt+"page user link failed")
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
			ID:        item.ID,
			UserID:    item.UserID,
			OriginURL: item.OriginURL,
			ShortURL:  item.ShortURL,
			CreatedAt: timex.FormatDateTime(item.CreatedAt),
			UpdatedAt: timex.FormatDateTime(item.UpdatedAt),
		}
		if item.ExpiredAt > 0 {
			slLink.ExpiredAt = timex.FormatDateTime(item.ExpiredAt)
		}
		data = append(data, slLink)
	}

	resp.Data = data
	resp.Total = total
	return resp, nil
}

func (svc *LinkService) DeleteLink(ctx context.Context, req *request.DelLinkReq) error {
	logFmt := "[LinkService][DeleteLink]"
	short, err := svc.linkRepo.GetByShort(ctx, req.ShortUrl)
	if err != nil {
		logs.Error(err, logFmt+"get by short link failed", zap.Any("shortLink", req.ShortUrl))
		return err
	}
	if short == nil {
		return errors.New("短链不存在")
	}
	if short.UserID != req.UserId {
		return errors.New("操作非法")
	}
	if err = svc.linkRepo.DeleteByShort(ctx, short.ShortURL); err != nil {
		logs.Error(err, logFmt+"delete by short failed")
		return err
	}
	gox.Run(context.Background(), func(ctx context.Context) {
		if err := bus.GetEventBus().Publish(ctx, consts.DeleteShortURLEvent, req.ShortUrl); err != nil {
			logs.Error(err, logFmt+fmt.Sprintf("publish %s event failed", consts.DeleteShortURLEvent))
		}
	})

	return err
}

func (svc *LinkService) validateOriginUrl(originUrl string) error {
	if originUrl == "" {
		return consts.ErrURLIsEmpty
	}
	parseUrl, err := url.Parse(originUrl)
	if err != nil {
		return err
	}
	if parseUrl.Scheme == "" {
		return consts.ErrSchemeIsEmpty
	}
	if parseUrl.Scheme != consts.HTTPScheme && parseUrl.Scheme != consts.HTTPSScheme {
		return consts.ErrSchemeInvalid
	}
	if parseUrl.Host == "" {
		return consts.ErrHostIsEmpty
	}
	return nil
}
