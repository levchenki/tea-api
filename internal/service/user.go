package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas/userSchemas"
	"sort"
	"strings"
	"time"
)

type UserRepository interface {
	Exists(telegramId uint64) (bool, error)
	Create(user *entity.User) error
	GetByTelegramId(telegramId uint64) (*entity.User, error)
	GetById(userId uuid.UUID) (*entity.User, error)
	SaveRefreshToken(userId, refreshTokenId uuid.UUID) error
	IsRefreshTokenExists(userId, refreshTokenId uuid.UUID) (bool, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) AuthenticateUser(tgUser *userSchemas.TelegramUser, botToken, jwtSecret string) (*userSchemas.UserTokens, error) {
	if err := s.verifyTelegramAuth(tgUser, botToken); err != nil {
		errResponse := errx.NewForbiddenError(fmt.Errorf("telegram verification error: %w", err))
		return nil, errResponse
	}

	exists, err := s.userRepository.Exists(tgUser.Id)
	if err != nil {
		errResponse := errx.NewInternalServerError(fmt.Errorf("user exists check error: %w", err))
		return nil, errResponse
	}

	if !exists {
		emptyUser := entity.NewEmptyUser(tgUser.Id, tgUser.FirstName, tgUser.LastName, tgUser.Username)
		err := s.userRepository.Create(emptyUser)
		if err != nil {
			errResponse := errx.NewInternalServerError(fmt.Errorf("user creation error: %w", err))
			return nil, errResponse
		}
	}

	u, err := s.userRepository.GetByTelegramId(tgUser.Id)

	if err != nil {
		errResponse := errx.NewInternalServerError(fmt.Errorf("user getting error: %w", err))
		return nil, errResponse
	}

	accessToken, err := s.generateAccessToken(u, jwtSecret)
	if err != nil {
		errResponse := errx.NewInternalServerError(fmt.Errorf("accessToken generation error: %w", err))
		return nil, errResponse
	}

	refreshToken, err := s.createRefreshToken(u.Id, jwtSecret)
	if err != nil {
		return nil, err
	}

	tokens := userSchemas.NewUserTokens(accessToken, refreshToken)
	return tokens, nil
}

func (s *UserService) createRefreshToken(userId uuid.UUID, jwtSecret string) (*userSchemas.RefreshToken, error) {
	refreshToken, err := s.generateRefreshToken(userId, jwtSecret)
	if err != nil {
		errResponse := errx.NewInternalServerError(fmt.Errorf("refreshToken generation error: %w", err))
		return nil, errResponse
	}

	err = s.userRepository.SaveRefreshToken(userId, refreshToken.Claims.Jit)
	if err != nil {
		return nil, err
	}
	return refreshToken, err
}

func (s *UserService) CheckAuthToken(authHeader string, jwtSecret string) (*userSchemas.AccessTokenClaims, error) {
	if !strings.HasPrefix(authHeader, "Bearer") {
		errResponse := errx.NewUnauthorizedError(fmt.Errorf("user is not authorized"))
		return nil, errResponse
	}

	accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
	accessToken, err := s.parseJWT(accessTokenString, jwtSecret)

	if err != nil {
		var errResponse *errx.AppError
		if errors.Is(err, jwt.ErrTokenExpired) {
			errResponse = errx.NewUnauthorizedError(fmt.Errorf("access token is expired"))
		} else {
			errResponse = errx.NewUnauthorizedError(err)
		}
		return nil, errResponse
	}

	accessTokenClaims, err := s.validateAccessToken(accessToken)

	if err != nil {
		errResponse := errx.NewUnauthorizedError(err)
		return nil, errResponse
	}

	return accessTokenClaims, nil
}

func (s *UserService) UpdateAccessToken(signedRefreshToken, jwtSecret string) (*userSchemas.UserTokens, error) {

	refreshToken, err := s.parseJWT(signedRefreshToken, jwtSecret)

	if err != nil {
		var errResponse *errx.AppError
		if errors.Is(err, jwt.ErrTokenExpired) {
			errResponse = errx.NewUnauthorizedError(fmt.Errorf("refresh token is expired"))
		} else {
			errResponse = errx.NewUnauthorizedError(err)
		}
		return nil, errResponse
	}

	rtClaims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		errResponse := errx.NewUnauthorizedError(err)
		return nil, errResponse
	}

	isExists, err := s.userRepository.IsRefreshTokenExists(rtClaims.Id, rtClaims.Jit)

	if err != nil {
		return nil, err
	}

	if !isExists {
		errResponse := errx.NewUnauthorizedError(fmt.Errorf("refresh token for user %s is not found", rtClaims.Id))
		return nil, errResponse
	}

	user, err := s.userRepository.GetById(rtClaims.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		errResponse := errx.NewNotFoundError(fmt.Errorf("user with id %s is not found", rtClaims.Id))
		return nil, errResponse
	}

	newAccessToken, err := s.generateAccessToken(user, jwtSecret)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.createRefreshToken(user.Id, jwtSecret)
	if err != nil {
		return nil, err
	}

	tokens := userSchemas.NewUserTokens(newAccessToken, newRefreshToken)
	return tokens, nil
}

func (s *UserService) verifyTelegramAuth(tgUser *userSchemas.TelegramUser, botToken string) error {
	if botToken == "" {
		return fmt.Errorf("invalid bot token")
	}

	var dataStrings []string
	dataStrings = append(dataStrings, fmt.Sprintf("auth_date=%d", tgUser.AuthDate))
	dataStrings = append(dataStrings, fmt.Sprintf("first_name=%s", tgUser.FirstName))
	dataStrings = append(dataStrings, fmt.Sprintf("id=%d", tgUser.Id))
	if tgUser.LastName != "" {
		dataStrings = append(dataStrings, fmt.Sprintf("last_name=%s", tgUser.LastName))
	}
	if tgUser.Username != "" {
		dataStrings = append(dataStrings, fmt.Sprintf("username=%s", tgUser.Username))
	}
	if tgUser.PhotoURL != "" {
		dataStrings = append(dataStrings, fmt.Sprintf("photo_url=%s", tgUser.PhotoURL))
	}

	sort.Strings(dataStrings)
	checkString := strings.Join(dataStrings, "\n")

	h := sha256.New()
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(checkString))
	hash := hex.EncodeToString(mac.Sum(nil))

	isValidHash := hash == tgUser.Hash
	if !isValidHash {
		return fmt.Errorf("hashes are not equal")
	}
	return nil
}

func (s *UserService) parseJWT(tokenString, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}

func (s *UserService) validateAccessToken(token *jwt.Token) (*userSchemas.AccessTokenClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userClaims := &userSchemas.AccessTokenClaims{}

	if exp, err := claims.GetExpirationTime(); err != nil || exp == nil {
		return nil, fmt.Errorf("exp can not be null")
	} else {
		userClaims.Exp = exp
	}

	if id, ok := claims["id"]; !ok {
		return nil, fmt.Errorf("invalid claims: id can not be null")
	} else {
		idUUID, err := uuid.Parse(id.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid claims: invalid id")
		}
		userClaims.Id = idUUID
	}

	if firstName, ok := claims["firstName"]; !ok {
		return nil, fmt.Errorf("invalid claims: firstName can not be null")
	} else {
		firstNameString, ok := firstName.(string)
		if ok {
			userClaims.FirstName = firstNameString
		} else {
			return nil, fmt.Errorf("invalid claims: invalid firstName")
		}
	}

	if username, ok := claims["username"]; !ok {
		return nil, fmt.Errorf("invalid claims: username can not be null")
	} else {
		usernameString, ok := username.(string)
		if ok {
			userClaims.Username = usernameString
		} else {
			return nil, fmt.Errorf("invalid claims: invalid username")
		}
	}

	if role, ok := claims["role"]; !ok {
		return nil, fmt.Errorf("invalid claims: role can not be null")
	} else {
		roleString, ok := role.(string)
		if ok {
			userClaims.Role = roleString
		} else {
			return nil, fmt.Errorf("invalid claims: invalid role")
		}
	}

	return userClaims, nil
}

func (s *UserService) validateRefreshToken(token *jwt.Token) (*userSchemas.RefreshTokenClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	refreshClaims := &userSchemas.RefreshTokenClaims{}

	if exp, err := claims.GetExpirationTime(); err != nil || exp == nil {
		return nil, fmt.Errorf("exp can not be null")
	} else {
		refreshClaims.Exp = exp
	}

	if id, ok := claims["id"]; !ok {
		return nil, fmt.Errorf("invalid claims: id can not be null")
	} else {
		idUUID, err := uuid.Parse(id.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid claims: invalid id")
		}
		refreshClaims.Id = idUUID
	}

	if jit, ok := claims["jit"]; !ok {
		return nil, fmt.Errorf("invalid claims: jit can not be null")
	} else {
		jitUUid, err := uuid.Parse(jit.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid claims: invalid jit")
		}
		refreshClaims.Jit = jitUUid
	}

	return refreshClaims, nil
}

func (s *UserService) generateAccessToken(user *entity.User, jwtSecret string) (*userSchemas.AccessToken, error) {
	expTime := time.Now().Add(time.Hour * 1)
	var role string
	if user.IsAdmin {
		role = "admin"
	} else {
		role = "user"
	}
	accessClaims := jwt.MapClaims{
		"id":        user.Id.String(),
		"firstName": user.FirstName,
		"username":  user.Username,
		"role":      role,
		"exp":       expTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedToken, err := token.SignedString([]byte(jwtSecret))

	accessToken := userSchemas.NewAccessToken(signedToken, user, role, expTime)
	return accessToken, err
}

func (s *UserService) generateRefreshToken(id uuid.UUID, jwtSecret string) (*userSchemas.RefreshToken, error) {
	expTime := time.Now().Add(time.Hour * 24 * 7)
	rtClaims := userSchemas.NewRefreshTokenClaims(id, expTime)

	claims := jwt.MapClaims{
		"id":  rtClaims.Id.String(),
		"exp": rtClaims.Exp.Unix(),
		"jit": rtClaims.Jit,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))

	refreshToken := userSchemas.NewRefreshToken(signedToken, rtClaims)

	return refreshToken, err
}
