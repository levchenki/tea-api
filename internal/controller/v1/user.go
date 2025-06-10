package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/logx"
	"github.com/levchenki/tea-api/internal/schemas/userSchemas"
	"net/http"
)

type UserService interface {
	AuthenticateUser(tgUser *userSchemas.TelegramUser, botToken, jwtSecret string) (*userSchemas.UserTokens, error)
	AuthenticateTelegramMiniApp(initData, botToken, jwtSecret string) (*userSchemas.UserTokens, error)
	CheckAuthToken(authHeader string, jwtSecret string) (*userSchemas.AccessTokenClaims, error)
	UpdateAccessToken(signedRefreshToken, jwtSecret string) (*userSchemas.UserTokens, error)
}

type UserController struct {
	jwtSecret   string
	botToken    string
	userService UserService
	log         logx.AppLogger
}

func NewUserController(jwtSecret, botToken string, userService UserService, log logx.AppLogger) *UserController {
	return &UserController{
		jwtSecret:   jwtSecret,
		botToken:    botToken,
		userService: userService,
		log:         log,
	}
}

// Auth godoc
//
//	@Summary	Auth
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		telegramUser	body		userSchemas.TelegramUser	true	"Telegram user"
//	@Success	200				{object}	userSchemas.TokenResponse
//	@Failure	400				{object}	errx.AppError
//	@Failure	403				{object}	errx.AppError
//	@Failure	500				{object}	errx.AppError
//	@Router		/api/v1/auth [post]
func (c *UserController) Auth(w http.ResponseWriter, r *http.Request) {
	var tgUser *userSchemas.TelegramUser
	if err := json.NewDecoder(r.Body).Decode(&tgUser); err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid data format: %w", err))
		handleError(w, r, c.log, errResponse)
		return
	}

	tokens, err := c.userService.AuthenticateUser(tgUser, c.botToken, c.jwtSecret)

	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken.SignedValue,
		Expires:  tokens.RefreshToken.Claims.Exp.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	render.JSON(w, r, tokens)
}

// AuthMiniApp godoc
//
//	@Summary	Auth
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		MiniAppInitData	body		userSchemas.MiniAppInitRequest	true	"Telegram init data"
//	@Success	200				{object}	userSchemas.TokenResponse
//	@Failure	400				{object}	errx.AppError
//	@Failure	403				{object}	errx.AppError
//	@Failure	500				{object}	errx.AppError
//	@Router		/api/v1/auth/mini-app [post]
func (c *UserController) AuthMiniApp(w http.ResponseWriter, r *http.Request) {
	var request userSchemas.MiniAppInitRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errResponse := errx.NewBadRequestError(fmt.Errorf("invalid data format: %w", err))
		handleError(w, r, c.log, errResponse)
		return
	}

	tokens, err := c.userService.AuthenticateTelegramMiniApp(request.InitData, c.botToken, c.jwtSecret)

	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken.SignedValue,
		Expires:  tokens.RefreshToken.Claims.Exp.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	render.JSON(w, r, tokens)
}

// UpdateAccessToken godoc
//
//	@Summary	Update access token
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		userId	path		string	true	"User ID"
//	@Success	200		{object}	userSchemas.TokenResponse
//	@Failure	400		{object}	errx.AppError
//	@Failure	403		{object}	errx.AppError
//	@Failure	500		{object}	errx.AppError
//	@Router		/api/v1/auth/refresh [post]
func (c *UserController) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	requestCookie, err := r.Cookie("refreshToken")

	if err != nil {
		var errResponse *errx.AppError
		if errors.Is(err, http.ErrNoCookie) {
			errResponse = errx.NewUnauthorizedError(fmt.Errorf("request has no refresh token"))
		} else {
			errResponse = errx.NewInternalServerError(err)
		}

		handleError(w, r, c.log, errResponse)
		return
	}

	tokens, err := c.userService.UpdateAccessToken(requestCookie.Value, c.jwtSecret)

	if err != nil {
		handleError(w, r, c.log, err)
		return
	}

	newCookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken.SignedValue,
		Expires:  tokens.RefreshToken.Claims.Exp.Time,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, newCookie)
	render.JSON(w, r, tokens)

}

func (c *UserController) AuthMiddleware(required bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !required && authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			accessTokenClaims, err := c.userService.CheckAuthToken(authHeader, c.jwtSecret)

			if err != nil {
				handleError(w, r, c.log, err)
				return
			}

			ctx := context.WithValue(r.Context(), "accessTokenClaims", accessTokenClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (c *UserController) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims := r.Context().Value("accessTokenClaims").(*userSchemas.AccessTokenClaims)

		if userClaims.Role == "admin" {
			next.ServeHTTP(w, r)
		} else {
			err := fmt.Errorf("user with id %s does not have admin permissions", userClaims.Id)
			errResponse := errx.NewForbiddenError(err)
			handleError(w, r, c.log, errResponse)
			return
		}
	})
}
