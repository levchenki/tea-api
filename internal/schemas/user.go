package schemas

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TelegramUser struct {
	Id        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
	PhotoURL  string `json:"photo_url"`
}

type UserClaims struct {
	Id        uuid.UUID       `json:"id"`
	FirstName string          `json:"firstName"`
	Username  string          `json:"username,omitempty"`
	Role      string          `json:"role"`
	Exp       jwt.NumericDate `json:"exp"`
}
