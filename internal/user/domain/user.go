package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) ValidateUserName() error {
	length := len(u.Username)
	if length >= 6 && length <= 30 {
		return nil
	}
	return errors.New("用户名长度不能小于6且不能大于30")
}
func (u *User) ValidatePassword() error {
	length := len(u.Password)
	if length < 8 || length > 64 {
		return errors.New("密码长度不能小于8或不能大于64")
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(u.Password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(u.Password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(u.Password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(u.Password)
	if !hasLower {
		return errors.New("密码中不包含小写字母")
	}
	if !hasUpper {
		return errors.New("密码中不包含大写字母")
	}
	if !hasDigit {
		return errors.New("密码中不包含数字")
	}
	if !hasSpecial {
		return errors.New("密码中不包含特殊字符")
	}

	return nil
}

func (u *User) EncryptPwd() (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	return string(hashPwd), err
}
