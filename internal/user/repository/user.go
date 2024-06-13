package repository

import (
	"context"
	"short-link/internal/user/repository/db"
)

type IUserRepository interface {
	AddUser(ctx context.Context, user *db.SlUser) error
	GetByUname(ctx context.Context, uname string) (*db.SlUser, error)
	GetByUnameAndPwd(ctx context.Context, uname, pwd string) (*db.SlUser, error)
}

type UserRepository struct {
}

func (repo *UserRepository) AddUser(ctx context.Context, user *db.SlUser) error {
	return db.NewSlUserDao(ctx).Create(user)
}

func (repo *UserRepository) GetByUname(ctx context.Context, uname string) (*db.SlUser, error) {
	return db.NewSlUserDao(ctx).GetByUname(uname)
}

func (repo *UserRepository) GetByUnameAndPwd(ctx context.Context, uname, pwd string) (*db.SlUser, error) {
	return db.NewSlUserDao(ctx).GetByUnameAndPwd(uname, pwd)
}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}
