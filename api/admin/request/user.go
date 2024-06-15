package request

// RegisterReq 注册请求参数
type RegisterReq struct {
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
