package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type AuthHandler struct {
	signUpUsecase                   usecasePorts.SignUpUsecase
	sendEmailVerificationUsecase    usecasePorts.SendEmailVerificationUsecase
	confirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase
}

func NewAuthHandler(signUpUsecase usecasePorts.SignUpUsecase, sendEmailVerificationUsecase usecasePorts.SendEmailVerificationUsecase, confirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase) *AuthHandler {
	return &AuthHandler{
		signUpUsecase:                   signUpUsecase,
		sendEmailVerificationUsecase:    sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase: confirmEmailVerificationUsecase,
	}
}

// SignUp godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя с указанным email и паролем
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body SignUpRequest true "Данные для регистрации"
// @Success      201 {object} SignUpResponse "Пользователь успешно зарегистрирован"
// @Failure      400 {object} ErrorResponse "Неверные данные запроса"
// @Failure      409 {object} ErrorResponse "Пользователь уже существует"
// @Failure      500 {object} ErrorResponse "Внутренняя ошибка сервера"
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
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body VerifyEmailRequest true "Email для отправки кода верификации"
// @Success      200 {object} VerifyEmailResponse "Email успешно отправлен"
// @Failure      400 {object} ErrorResponse "Неверный формат данных"
// @Failure      404 {object} ErrorResponse "Пользователь не найден"
// @Failure      500 {object} ErrorResponse "Внутренняя ошибка сервера"
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
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body ConfirmEmailRequest true "Email для отправки кода верификации"
// @Success      200 {object} ConfirmEmailResponse "Email успешно отправлен"
// @Failure      400 {object} ErrorResponse "Неверный формат данных"
// @Failure      404 {object} ErrorResponse "Пользователь не найден"
// @Failure      500 {object} ErrorResponse "Внутренняя ошибка сервера"
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
