package userSchemas

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"time"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type AccessTokenClaims struct {
	Id        uuid.UUID        `json:"id"`
	FirstName string           `json:"firstName"`
	Username  string           `json:"username,omitempty"`
	Role      string           `json:"role"`
	Exp       *jwt.NumericDate `json:"exp"`
}

type RefreshTokenClaims struct {
	Id  uuid.UUID        `json:"id"`
	Exp *jwt.NumericDate `json:"exp"`
	Jit uuid.UUID        `json:"jit"`
}

func NewRefreshTokenClaims(userId uuid.UUID, expTime time.Time) *RefreshTokenClaims {
	return &RefreshTokenClaims{
		Id:  userId,
		Exp: jwt.NewNumericDate(expTime),
		Jit: uuid.New(),
	}
}

type AccessToken struct {
	SignedValue string             `json:"_"`
	Claims      *AccessTokenClaims `json:"-"`
}
type RefreshToken struct {
	SignedValue string              `json:"-"`
	Claims      *RefreshTokenClaims `json:"-"`
}

type UserTokens struct {
	AccessToken  *AccessToken  `json:"_"`
	RefreshToken *RefreshToken `json:"-"`
}

func NewAccessToken(signedValue string, user *entity.User, role string, expTime time.Time) *AccessToken {
	return &AccessToken{
		SignedValue: signedValue,
		Claims: &AccessTokenClaims{
			Id:        user.Id,
			FirstName: user.FirstName,
			Username:  user.Username,
			Role:      role,
			Exp:       jwt.NewNumericDate(expTime),
		},
	}
}

func NewRefreshToken(signedValue string, claims *RefreshTokenClaims) *RefreshToken {
	return &RefreshToken{
		SignedValue: signedValue,
		Claims:      claims,
	}
}

func NewUserTokens(accessToken *AccessToken, refreshToken *RefreshToken) *UserTokens {
	userTokens := &UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return userTokens
}

func (ut *UserTokens) MarshalJSON() ([]byte, error) {
	type tokenOnly struct {
		// todo token -> accessToken
		Token string `json:"token"`
	}
	return json.Marshal(&tokenOnly{
		Token: ut.AccessToken.SignedValue,
	})
}
