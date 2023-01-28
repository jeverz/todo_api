package session

type Session struct {
	Id       int64  `json:"i"`
	Timeout  int64  `json:"t"`
	RefreshTimeout int64 `json:"r"`
	IsAdmin  bool  `json:"a"`
}
