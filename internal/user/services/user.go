package services

import (
	"context"
	"errors"
	"short-link/api/admin/request"
	"short-link/internal/user/domain"
	"short-link/internal/user/repository"
	"short-link/internal/user/repository/db"
	"short-link/logs"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(),
	}
}

func (svc *UserService) Register(ctx context.Context, req *request.RegisterReq) error {
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
	if u != nil {
		return errors.New("用户名已存在")
	}
	pwd, err := user.EncryptPwd()
	if err != nil {
		return err
	}
	// 密码加密
	err = svc.userRepository.AddUser(ctx, &db.SlUser{
		Username: req.Username,
		Password: pwd,
	})

	return err
}

func (svc *UserService) Login(ctx context.Context, req *request.LoginReq) (*db.SlUser, error) {
	userModel, err := svc.userRepository.GetByUname(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if userModel == nil {
		logs.Error(err, "[UserService][Login] user not found", zap.Any("username", req.Username))
		return nil, errors.New("user not found")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.Password)); err != nil {
		logs.Error(err, "[UserService][Login] compare password failed",
			zap.Any("reqPassword", req.Password),
			zap.Any("password", userModel.Password))
		return nil, errors.New("wrong password")
	}

	return userModel, err
}
