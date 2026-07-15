package api

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 返回给前端的用户信息
type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// LoginResponse 登录成功后的返回结果
type LoginResponse struct {
	Token string       `json:"token"` // 发放给前端的 JWT
	User  UserResponse `json:"user"`  // 用户基本信息（复用之前的 UserResponse）
}
