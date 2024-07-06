package services

import (
	"context"
	"errors"
	"fmt"
	"short-link/api/admin/request"
	"short-link/api/admin/resopnse"
	"short-link/ctxkit"
	"short-link/internal/link/repository"
	"short-link/logs"
	"short-link/utils/netx"
	"short-link/utils/timex"

	"go.uber.org/zap"
)

type BlackListService struct {
	linkRepo      repository.ILinkRepository
	blackListRepo repository.IBlackListRepository
}

func NewBlackListService() *BlackListService {
	return &BlackListService{
		linkRepo:      repository.NewLinkRepository(),
		blackListRepo: repository.NewBlackListRepository(),
	}
}

func (svc *BlackListService) AddBlackList(ctx context.Context, req *request.AddBlackListReq) error {
	logFmt := "[BlackListService][AddBlackList]"
	short, err := svc.linkRepo.GetByShort(ctx, req.ShortURL)
	if err != nil {
		logs.Error(err, logFmt+"get by short link failed", zap.Any("shortLink", req.ShortURL))
		return err
	}
	if short == nil {
		return errors.New("短链不存在")
	}
	if short.UserID != req.UserID {
		return errors.New("操作非法")
	}
	IPTint, err := netx.IPToInt(req.IP)
	if err != nil {
		logs.Error(err, logFmt+"transform IP to int failed")
		return err
	}
	// 添加黑名单
	if err = svc.blackListRepo.AddBlackList(ctx, req.ShortURL, IPTint); err != nil {
		logs.Error(err, logFmt+"add short url IP black list failed")
		return err
	}
	return nil
}

func (svc *BlackListService) DeleteBlackList(ctx context.Context, id uint64) error {
	logFmt := "[BlackListService][DeleteBlackList]"

	userId := ctxkit.GetUserId(ctx)
	blackList, err := svc.blackListRepo.GetByID(ctx, id)
	if err != nil {
		logs.Error(err, logFmt+"get black list by id failed")
		return err
	}
	if blackList == nil {
		err = fmt.Errorf("%d 黑名单不存", id)
		logs.Error(err, logFmt+"black list not found")
		return err
	}

	su := blackList.ShortURL

	shortUrl, err := svc.linkRepo.GetByShort(ctx, su)
	if err != nil {
		logs.Error(err, logFmt+"get short url failed")
		return err
	}
	if shortUrl == nil {
		err = errors.New("短链不存在")
		logs.Error(err, logFmt+"short url not found")
		return err
	}

	if shortUrl.UserID != userId {
		err = errors.New("操作非法")
		logs.Error(err, logFmt+"")
		return err
	}

	err = svc.blackListRepo.Delete(ctx, id)
	if err != nil {
		logs.Error(err, logFmt+"delete failed")
	}
	return err
}

func (svc *BlackListService) ListBlackList(ctx context.Context, req *request.ListBlackListReq) (*resopnse.ListBlackListResp, error) {
	logFmt := "[BlackListService][DeleteBlackList]"

	IP, err := netx.IPToInt(req.IP)
	if err != nil {
		return nil, err
	}
	blackList, total, err := svc.blackListRepo.PageBlackList(ctx, req.ShortUrl, IP, req.Page, req.PageSize)
	if err != nil {
		logs.Error(err, logFmt+"list black list failed")
		return nil, err
	}

	data := make([]resopnse.BlackList, 0, len(blackList))
	for _, item := range blackList {
		IPStr := netx.IntToIP(item.IP)
		data = append(data, resopnse.BlackList{
			ID:        item.Id,
			ShortUrl:  item.ShortURL,
			IP:        IPStr,
			CreatedAt: timex.FormatDateTime(item.CreatedAt),
		})
	}

	return &resopnse.ListBlackListResp{
		Data:  data,
		Total: total,
	}, nil
}
