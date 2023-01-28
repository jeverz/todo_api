package todoapi

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	IsAdmin  bool   `json:"isadmin"`
}

type RegisterRequest struct {
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
}
