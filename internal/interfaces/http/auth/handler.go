package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/auth_usecases"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type AuthHandler struct {
	signUpUsecase                   usecasePorts.SignUpUsecase
	sendEmailVerificationUsecase    usecasePorts.SendEmailVerificationUsecase
	confirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase
	signInUsecase                   usecasePorts.SignInUsecase
	refreshTokenUsecase             usecasePorts.RefreshTokenUsecase
}

func NewAuthHandler(
	signUpUsecase usecasePorts.SignUpUsecase,
	sendEmailVerificationUsecase usecasePorts.SendEmailVerificationUsecase,
	confirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase,
	signInUsecase usecasePorts.SignInUsecase,
	refreshTokenUsecase usecasePorts.RefreshTokenUsecase,
) *AuthHandler {
	return &AuthHandler{
		signUpUsecase:                   signUpUsecase,
		sendEmailVerificationUsecase:    sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase: confirmEmailVerificationUsecase,
		signInUsecase:                   signInUsecase,
		refreshTokenUsecase:             refreshTokenUsecase,
	}
}

// SignUp godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя с указанным email и паролем
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body SignUpRequest true "Данные для регистрации"
// @Success      201 {object} SignUpResponse "Пользователь успешно зарегистрирован"
// @Router       /v1/auth/sign-up [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	user, err := h.signUpUsecase.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response := ToSignUpResponse(user)

	c.JSON(http.StatusCreated, response)
}

// VerifyEmail godoc
// @Summary      Отправка кода верификации на email
// @Description  Отправляет код верификации на указанный email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body VerifyEmailRequest true "Email для отправки кода верификации"
// @Success      200 {object} VerifyEmailResponse "Email успешно отправлен"
// @Router       /v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	err = h.sendEmailVerificationUsecase.Execute(c.Request.Context(), req.Email)

	if err != nil {
		_ = c.Error(err)
		return
	}

	response := VerifyEmailResponse{
		Message: "Email успешно отправлен",
	}
	c.JSON(http.StatusOK, response)
}

// ConfirmEmail  godoc
// @Summary      Подтверждение кода верификации на email
// @Description  Подтверждение кода
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body ConfirmEmailRequest true "Email для отправки кода верификации"
// @Success      200 {object} ConfirmEmailResponse "Email успешно отправлен"
// @Router       /v1/auth/confirm-email [post]
func (h *AuthHandler) ConfirmEmail(c *gin.Context) {
	var req ConfirmEmailRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	err = h.confirmEmailVerificationUsecase.Execute(c.Request.Context(), req.Email, req.Code)

	if err != nil {
		_ = c.Error(err)
		return
	}

	response := ConfirmEmailResponse{
		Message: "Email успешно подтвержден",
	}
	c.JSON(http.StatusOK, response)
}

// SignIn godoc
// @Summary      Вход пользователя
// @Description  Аутентификация пользователя с помощью email и пароля
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body SignInRequest true "Данные для входа"
// @Success      200 {object} SignInResponse "Пользователь успешно аутентифицирован"
// @Router       /v1/auth/sign-in [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	user, tokens, err := h.signInUsecase.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := SignInResponse{
		User:   ToSignUpResponse(user),
		Tokens: *tokens,
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary      Обновление access токена
// @Description  Обновляет access токен с помощью refresh токена
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Данные для обновления токена"
// @Success      200 {object} RefreshResponse "Токен успешно обновлен"
// @Router       /v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	newAccess, err := h.refreshTokenUsecase.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := RefreshResponse{AccessToken: newAccess}
	c.JSON(http.StatusOK, resp)
}
