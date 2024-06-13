package services

import (
	"context"
	"errors"
	"short-link/api/request"
)

type LinkService struct {
}

func NewLinkService() {

}

func (svc *LinkService) AddLink(ctx context.Context, req *request.AddLinkReq) error {
	if req.OriginalUrl == "" {
		return errors.New("原链接为空")
	}

	return nil
}
