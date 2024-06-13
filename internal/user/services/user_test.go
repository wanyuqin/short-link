package services

import (
	"context"
	"short-link/api/request"
	"short-link/internal/user/repository"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	type fields struct {
		userRepository repository.IUserRepository
	}
	type args struct {
		ctx context.Context
		req *request.Register
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "01",
			fields: fields{
				userRepository: repository.NewUserRepository(),
			},
			args: args{
				ctx: context.Background(),
				req: &request.Register{
					Username: "admin",
					Password: "Admin12345@",
				},
			},
			wantErr: true,
		},
		{
			name: "02",
			fields: fields{
				userRepository: repository.NewUserRepository(),
			},
			args: args{
				ctx: context.Background(),
				req: &request.Register{
					Username: "admin-1",
					Password: "Admin12345@",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &UserService{
				userRepository: tt.fields.userRepository,
			}
			if err := svc.Register(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
