package schemas

import "github.com/golang-jwt/jwt/v5"

type TelegramUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
	PhotoURL  string `json:"photo_url"`
}

type TelegramUserClaims struct {
	Id        int             `json:"id"`
	FirstName string          `json:"firstName"`
	Username  string          `json:"username,omitempty"`
	Exp       jwt.NumericDate `json:"exp"`
}
