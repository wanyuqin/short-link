package db

import (
	"context"
	"reflect"
	"testing"
)

func TestSlUser_Create(t *testing.T) {
	type fields struct {
		ID        int64
		Username  string
		Password  string
		Salt      string
		CreatedAt int64
		UpdatedAt int64
		IsDel     int
	}
	type args struct {
		ctx context.Context
		u   SlUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				u: SlUser{
					Username: "admin",
					Password: "admin",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlUser{
				ID:        tt.fields.ID,
				Username:  tt.fields.Username,
				Password:  tt.fields.Password,
				Salt:      tt.fields.Salt,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
				IsDel:     tt.fields.IsDel,
			}
			if err := m.Create(tt.args.ctx, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlUser_GetByUname(t *testing.T) {
	type fields struct {
		ID        int64
		Username  string
		Password  string
		Salt      string
		CreatedAt int64
		UpdatedAt int64
		IsDel     int
	}
	type args struct {
		ctx   context.Context
		uname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SlUser
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:   context.Background(),
				uname: "admin",
			},
			wantErr: false,
			want: SlUser{
				ID:        1,
				Username:  "admin",
				Password:  "admin",
				CreatedAt: 1718207749353,
				UpdatedAt: 1718207749353,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlUser{
				ID:        tt.fields.ID,
				Username:  tt.fields.Username,
				Password:  tt.fields.Password,
				Salt:      tt.fields.Salt,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
				IsDel:     tt.fields.IsDel,
			}
			got, err := m.GetByUname(tt.args.ctx, tt.args.uname)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByUname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByUname() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlUser_GetByUnameAndPwd(t *testing.T) {
	type fields struct {
		ID        int64
		Username  string
		Password  string
		Salt      string
		CreatedAt int64
		UpdatedAt int64
		IsDel     int
	}
	type args struct {
		ctx   context.Context
		uname string
		pwd   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SlUser
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:   context.Background(),
				uname: "admin",
				pwd:   "admin",
			},
			wantErr: false,
			want: SlUser{
				ID:        1,
				Username:  "admin",
				Password:  "admin",
				CreatedAt: 1718207749353,
				UpdatedAt: 1718207749353,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlUser{
				ID:        tt.fields.ID,
				Username:  tt.fields.Username,
				Password:  tt.fields.Password,
				Salt:      tt.fields.Salt,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
				IsDel:     tt.fields.IsDel,
			}
			got, err := m.GetByUnameAndPwd(tt.args.ctx, tt.args.uname, tt.args.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByUnameAndPwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByUnameAndPwd() got = %v, want %v", got, tt.want)
			}
		})
	}
}
