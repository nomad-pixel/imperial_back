package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type AuthHandler struct {
	signUpUsecase usecasePorts.SignUpUsecase
}

func NewAuthHandler(signUpUsecase usecasePorts.SignUpUsecase) *AuthHandler {
	return &AuthHandler{signUpUsecase: signUpUsecase}
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

	// Явное преобразование Entity -> DTO
	response := ToSignUpResponse(user)

	c.JSON(http.StatusCreated, response)
}
