package v1

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/levchenki/tea-api/internal/entity"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas"
	"net/http"
	"sort"
	"strings"
	"time"
)

type UserService interface {
	Create(user *entity.User) error
	Exists(telegramId uint64) (bool, error)
	GetByTelegramId(telegramId uint64) (*entity.User, error)
}

type UserController struct {
	jwtSecret   string
	botToken    string
	userService UserService
}

func NewUserController(jwtSecret, botToken string, userService UserService) *UserController {
	return &UserController{jwtSecret, botToken, userService}
}

func (c *UserController) Auth(w http.ResponseWriter, r *http.Request) {
	var telegramUser schemas.TelegramUser
	if err := json.NewDecoder(r.Body).Decode(&telegramUser); err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid data format: %w", err))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	if err := c.verifyTelegramAuth(telegramUser); err != nil {
		errResponse := errx.ErrorForbidden(fmt.Errorf("telegram verification error: %w", err))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	exists, err := c.userService.Exists(telegramUser.Id)
	if err != nil {
		errResponse := errx.ErrorInternalServer(fmt.Errorf("user exists check error: %w", err))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	if !exists {
		emptyUser := entity.NewEmptyUser(telegramUser.Id, telegramUser.FirstName, telegramUser.LastName, telegramUser.Username)
		err := c.userService.Create(emptyUser)
		if err != nil {
			errResponse := errx.ErrorInternalServer(fmt.Errorf("user creation error: %w", err))
			render.Status(r, errResponse.HTTPStatusCode)
			render.JSON(w, r, errResponse)
			return
		}
	}

	u, err := c.userService.GetByTelegramId(telegramUser.Id)

	if err != nil {
		errResponse := errx.ErrorInternalServer(fmt.Errorf("user getting error: %w", err))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	token, err := c.generateJWT(u, c.jwtSecret)
	if err != nil {
		errResponse := errx.ErrorInternalServer(fmt.Errorf("token generation error: %w", err))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	render.JSON(w, r, map[string]string{
		"token": token,
	})
}

func (c *UserController) verifyTelegramAuth(data schemas.TelegramUser) error {
	if c.botToken == "" {
		return fmt.Errorf("invalid bot token")
	}

	var dataStrings []string
	dataStrings = append(dataStrings, fmt.Sprintf("auth_date=%d", data.AuthDate))
	dataStrings = append(dataStrings, fmt.Sprintf("first_name=%s", data.FirstName))
	dataStrings = append(dataStrings, fmt.Sprintf("id=%d", data.Id))
	if data.LastName != "" {
		dataStrings = append(dataStrings, fmt.Sprintf("last_name=%s", data.LastName))
	}
	if data.Username != "" {
		dataStrings = append(dataStrings, fmt.Sprintf("username=%s", data.Username))
	}
	if data.PhotoURL != "" {
		dataStrings = append(dataStrings, fmt.Sprintf("photo_url=%s", data.PhotoURL))
	}

	sort.Strings(dataStrings)
	checkString := strings.Join(dataStrings, "\n")

	h := sha256.New()
	h.Write([]byte(c.botToken))
	secretKey := h.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(checkString))
	hash := hex.EncodeToString(mac.Sum(nil))

	isValidHash := hash == data.Hash
	if !isValidHash {
		return fmt.Errorf("hashes are not equal")
	}
	return nil
}

func (c *UserController) generateJWT(user *entity.User, jwtSecret string) (string, error) {
	expDate := time.Now().Add(time.Hour * 1).Unix()
	var role string
	if user.IsAdmin {
		role = "admin"
	} else {
		role = "user"
	}
	claims := jwt.MapClaims{
		"id":        user.Id.String(),
		"firstName": user.FirstName,
		"username":  user.Username,
		"role":      role,
		"exp":       expDate,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (c *UserController) AuthMiddleware(required bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !required && authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer") {
				errResponse := errx.ErrorUnauthorized(fmt.Errorf("user is not authorized"))
				render.Status(r, errResponse.HTTPStatusCode)
				render.JSON(w, r, errResponse)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := c.parseToken(tokenString, c.jwtSecret)

			if err != nil {
				errResponse := errx.ErrorUnauthorized(err)
				render.Status(r, errResponse.HTTPStatusCode)
				render.JSON(w, r, errResponse)
				return
			}

			claims, err := c.parseClaims(token)

			if err != nil {
				errResponse := errx.ErrorUnauthorized(err)
				render.Status(r, errResponse.HTTPStatusCode)
				render.JSON(w, r, errResponse)
				return
			}

			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (c *UserController) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims := r.Context().Value("user").(*schemas.UserClaims)

		if userClaims.Role == "admin" {
			next.ServeHTTP(w, r)
		} else {
			errResponse := errx.ErrorForbidden(fmt.Errorf("user does not have admin permissions"))
			render.Status(r, errResponse.HTTPStatusCode)
			render.JSON(w, r, errResponse)
			return
		}
	})
}

func (c *UserController) parseToken(tokenString, jwtSecret string) (*jwt.Token, error) {
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

func (c *UserController) parseClaims(token *jwt.Token) (*schemas.UserClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userClaims := schemas.UserClaims{}

	if exp, err := claims.GetExpirationTime(); err != nil || exp == nil {
		return nil, fmt.Errorf("exp can not be null")
	} else {
		userClaims.Exp = *exp
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

	return &userClaims, nil
}
