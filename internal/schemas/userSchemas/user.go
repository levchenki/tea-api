package userSchemas

type TelegramUser struct {
	Id        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
	PhotoURL  string `json:"photo_url"`
}

type MiniAppInitRequest struct {
	InitData string `json:"initData"`
}
