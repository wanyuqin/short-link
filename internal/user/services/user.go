package services

import (
	"context"
	"errors"
	"short-link/api/request"
	"short-link/internal/user/domain"
	"short-link/internal/user/repository"
	"short-link/internal/user/repository/db"
)

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(),
	}
}

func (svc *UserService) Register(ctx context.Context, req *request.Register) error {
	user := domain.User{
		Username: req.Username,
		Password: req.Password,
	}
	if err := user.ValidateUserName(); err != nil {
		return err
	}
	if err := user.ValidatePassword(); err != nil {
		return err
	}

	u, err := svc.userRepository.GetByUname(ctx, req.Username)
	if err != nil {
		return err
	}
	if u.ID > 0 {
		return errors.New("用户名已存在")
	}
	pwd, err := user.EncryptPwd()
	if err != nil {
		return err
	}
	// 密码加密
	err = svc.userRepository.AddUser(ctx, db.SlUser{
		Username: req.Username,
		Password: pwd,
	})

	return err
}
