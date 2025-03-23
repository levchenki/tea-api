package controller

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/schemas"
	"net/http"
	"sort"
	"strings"
	"time"
)

type AuthController struct {
	jwtSecret string
	botToken  string
}

func NewAuthController(jwtSecret, botToken string) *AuthController {
	return &AuthController{jwtSecret, botToken}
}

func (c *AuthController) Auth(w http.ResponseWriter, r *http.Request) {
	var user schemas.TelegramUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errResponse := errx.ErrorBadRequest(fmt.Errorf("invalid data format: %w", err))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	if err := c.verifyTelegramAuth(user); err != nil {
		errResponse := errx.ErrorForbidden(fmt.Errorf("telegram verification error: %w", err))
		render.Status(r, errResponse.HTTPStatusCode)
		render.JSON(w, r, errResponse)
		return
	}

	token, err := c.generateJWT(user, c.jwtSecret)
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

func (c *AuthController) verifyTelegramAuth(data schemas.TelegramUser) error {
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

func (c *AuthController) generateJWT(user schemas.TelegramUser, jwtSecret string) (string, error) {
	expDate := time.Now().Add(time.Hour * 1).Unix()
	claims := jwt.MapClaims{
		"id":        user.Id,
		"firstName": user.FirstName,
		"username":  user.Username,
		"exp":       expDate,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (c *AuthController) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
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

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims)
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

func (c *AuthController) parseToken(tokenString, jwtSecret string) (*jwt.Token, error) {
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

func (c *AuthController) parseClaims(token *jwt.Token) (*schemas.TelegramUserClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	var tgClaims schemas.TelegramUserClaims

	if exp, err := claims.GetExpirationTime(); err != nil || exp == nil {
		return nil, fmt.Errorf("exp can not be null")
	} else {
		tgClaims.Exp = *exp
	}

	if id, ok := claims["id"]; !ok {
		return nil, fmt.Errorf("invalid claims: id can not be null")
	} else {
		idInt, ok := id.(float64)
		if ok {
			tgClaims.Id = int(idInt)
		} else {
			return nil, fmt.Errorf("invalid claims: invalid id")
		}

	}

	if firstName, ok := claims["firstName"]; !ok {
		return nil, fmt.Errorf("invalid claims: firstName can not be null")
	} else {
		firstNameString, ok := firstName.(string)
		if ok {
			tgClaims.FirstName = firstNameString
		} else {
			return nil, fmt.Errorf("invalid claims: invalid firstName")
		}
	}

	if username, ok := claims["username"]; !ok {
		return nil, fmt.Errorf("invalid claims: username can not be null")
	} else {
		usernameString, ok := username.(string)
		if ok {
			tgClaims.Username = usernameString
		} else {
			return nil, fmt.Errorf("invalid claims: invalid username")
		}
	}

	return &tgClaims, nil
}
