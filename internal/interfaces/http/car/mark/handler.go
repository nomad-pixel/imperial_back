package car

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/car"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CarMarkHandler struct {
	createCarMark usecasePorts.CreateCarMarkUsecase
	getCarMark    usecasePorts.GetCarMarkUsecase
	getCarMarks   usecasePorts.GetCarMarksListUsecase
	updateCarMark usecasePorts.UpdateCarMarkUsecase
	deleteCarMark usecasePorts.DeleteCarMarkUsecase
}

func NewCarMarkHandler(
	createCarMark usecasePorts.CreateCarMarkUsecase,
	getCarMark usecasePorts.GetCarMarkUsecase,
	getCarMarks usecasePorts.GetCarMarksListUsecase,
	updateCarMark usecasePorts.UpdateCarMarkUsecase,
	deleteCarMark usecasePorts.DeleteCarMarkUsecase,
) *CarMarkHandler {
	return &CarMarkHandler{
		createCarMark: createCarMark,
		getCarMark:    getCarMark,
		getCarMarks:   getCarMarks,
		updateCarMark: updateCarMark,
		deleteCarMark: deleteCarMark,
	}
}

// CreateCarMark godoc
// @Summary      Создание новой марки машины
// @Description  Создает новую марку с указанным названием
// @Tags         Car Marks
// @Accept       json
// @Produce      json
// @Param        request body CreateCarMarkRequest true "Данные для создания марки"
// @Success      201 {object}  CarMarkResponse  "Марка успешно создана"
// @Security     BearerAuth
// @Router       /v1/cars/car-marks [post]
func (h *CarMarkHandler) CreateCarMark(c *gin.Context) {
	var req CreateCarMarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	mark, err := h.createCarMark.Execute(c.Request.Context(), req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, mark)
}

// GetCarMark godoc
// @Summary      Получение марки по ID
// @Description  Возвращает информацию о марке по указанному ID
// @Tags         Car Marks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID марки"
// @Success      200 {object}  CarMarkResponse  "Информация о марке"
// @Security     BearerAuth
// @Router       /v1/cars/car-marks/{id} [get]
func (h *CarMarkHandler) GetCarMark(c *gin.Context) {
	markID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID марки"))
		return
	}

	mark, err := h.getCarMark.Execute(c.Request.Context(), markID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mark)
}

// GetCarMarks godoc
// @Summary      Получение списка марок
// @Description  Возвращает список всех марок с поддержкой пагинации
// @Tags         Car Marks
// @Accept       json
// @Produce      json
// @Param        offset query int false "Смещение для пагинации" default(0)
// @Param        limit query int false "Лимит для пагинации" default(20)
// @Success      200 {object}  ListCarMarksResponse  "Список марок"
// @Security     BearerAuth
// @Router       /v1/cars/car-marks [get]
func (h *CarMarkHandler) GetCarMarks(c *gin.Context) {
	offset := int64(0)
	limit := int64(20)

	if o, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		offset = o
	}
	if l, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		limit = l
	}

	total, marks, err := h.getCarMarks.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ListCarMarksResponse{
		Total: total,
		Data:  marks,
	})
}

// UpdateCarMark godoc
// @Summary      Обновление марки
// @Description  Обновляет информацию о марке по указанному ID
// @Tags         Car Marks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID марки"
// @Param        request body UpdateCarMarkRequest true "Новые данные для марки"
// @Success      200 {object}  CarMarkResponse  "Марка успешно обновлена"
// @Security     BearerAuth
// @Router       /v1/cars/car-marks/{id} [put]
func (h *CarMarkHandler) UpdateCarMark(c *gin.Context) {
	markID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID марки"))
		return
	}

	var req UpdateCarMarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	mark, err := h.updateCarMark.Execute(c.Request.Context(), markID, req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mark)
}

// DeleteCarMark godoc
// @Summary      Удаление марки
// @Description  Удаляет марку по указанному ID
// @Tags         Car Marks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID марки"
// @Success      200 {object}  MessageResponse  "Марка успешно удалена"
// @Security     BearerAuth
// @Router       /v1/cars/car-marks/{id} [delete]
func (h *CarMarkHandler) DeleteCarMark(c *gin.Context) {
	markID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID марки"))
		return
	}

	err = h.deleteCarMark.Execute(c.Request.Context(), markID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Марка успешно удалена"})
}
