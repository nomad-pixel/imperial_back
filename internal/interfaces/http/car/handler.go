package car

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CarHandler struct {
	createCar usecasePorts.CreateCarUsecase
}

func NewCarHandler(
	createCar usecasePorts.CreateCarUsecase,
) *CarHandler {
	return &CarHandler{
		createCar: createCar,
	}
}

// CreateCar godoc
// @Summary      Создание нового автомобиля
// @Description  Создает новый автомобиль с указанными параметрами
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        request body CreateCarRequest true "Данные для создания автомобиля"
// @Success      200 {object}  CarResponse  "Автомобиль успешно создан"
// @Security     BearerAuth
// @Router       /v1/cars [post]
func (h *CarHandler) CreateCar(c *gin.Context) {
	var req CreateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	car := &entities.Car{
		Name:           req.Name,
		ImageUrl:       req.ImageUrl,
		OnlyWithDriver: req.OnlyWithDriver,
		PricePerDay:    req.PricePerDay,
		Mark: &entities.CarMark{
			ID: req.MarkId,
		},
		Category: &entities.CarCategory{
			ID: req.CategoryId,
		},

		Tags: make([]*entities.CarTag, 0, len(req.TagsIds)),
	}

	for _, tagID := range req.TagsIds {
		car.Tags = append(car.Tags, &entities.CarTag{ID: tagID})
	}

	createdCar, err := h.createCar.Execute(c.Request.Context(), car)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, NewCarResponse(createdCar))
}
