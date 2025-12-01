package celebrity

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/celebrity"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CelebrityHandler struct {
	createCelebrityUsecase usecasePorts.CreateCelebrityUsecase
}

func NewCelebrityHandler(
	createCelebrityUsecase usecasePorts.CreateCelebrityUsecase,
) *CelebrityHandler {
	return &CelebrityHandler{
		createCelebrityUsecase: createCelebrityUsecase,
	}
}

// CreateCelebrity godoc
// @Summary Create a new celebrity
// @Description Create a new celebrity with the provided details
// @Tags Celebrities
// @Accept json
// @Produce json
// @Param celebrity body CreateCelebrityRequest true "Celebrity data"
// @Success 201 {object} entities.Celebrity
// @Router /v1/celebrities [post]
// @Security     BearerAuth
func (h *CelebrityHandler) CreateCelebrity(c *gin.Context) {
	var req CreateCelebrityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	celebrity := &entities.Celebrity{
		Name:  req.Name,
		Image: req.Image,
	}

	err := h.createCelebrityUsecase.Execute(c.Request.Context(), celebrity)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(201, celebrity)
}
