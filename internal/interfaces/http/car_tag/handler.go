package car_tag

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/car_usecases"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CarTagHandler struct {
	createCarTag usecasePorts.CreateCarTagUsecase
	getCarTag    usecasePorts.GetCarTagUsecase
	getCarTags   usecasePorts.GetCarTagsListUsecase
	updateCarTag usecasePorts.UpdateCarTagUsecase
	deleteCarTag usecasePorts.DeleteCarTagUsecase
}

func NewCarTagHandler(
	createCarTag usecasePorts.CreateCarTagUsecase,
	getCarTag usecasePorts.GetCarTagUsecase,
	getCarTags usecasePorts.GetCarTagsListUsecase,
	updateCarTag usecasePorts.UpdateCarTagUsecase,
	deleteCarTag usecasePorts.DeleteCarTagUsecase,
) *CarTagHandler {
	return &CarTagHandler{
		createCarTag: createCarTag,
		getCarTag:    getCarTag,
		getCarTags:   getCarTags,
		updateCarTag: updateCarTag,
		deleteCarTag: deleteCarTag,
	}
}

// CreateCarTag godoc
// @Summary      Создание нового тега для машины
// @Description  Создает новый тег с указанным названием
// @Tags         Car Tags
// @Accept       json
// @Produce      json
// @Param        request body CreateCarTagRequest true "Данные для создания тега"
// @Success      201 {object}  CarTagResponse  "Тег успешно создан"
// @Security     BearerAuth
// @Router       /v1/cars/car-tags [post]
func (h *CarTagHandler) CreateCarTag(c *gin.Context) {
	var req CreateCarTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	tag, err := h.createCarTag.Execute(c.Request.Context(), req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetCarTag godoc
// @Summary      Получение тега по ID
// @Description  Возвращает информацию о теге по указанному ID
// @Tags         Car Tags
// @Accept       json
// @Produce      json
// @Param        id path int true "ID тега"
// @Success      200 {object}  CarTagResponse  "Информация о теге"
// @Router       /v1/cars/car-tags/{id} [get]
func (h *CarTagHandler) GetCarTag(c *gin.Context) {
	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID тега"))
		return
	}

	tag, err := h.getCarTag.Execute(c.Request.Context(), tagID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tag)
}

// GetCarTags godoc
// @Summary      Получение списка тегов
// @Description  Возвращает список всех тегов с поддержкой пагинации
// @Tags         Car Tags
// @Accept       json
// @Produce      json
// @Param        offset query int false "Смещение для пагинации" default(0)
// @Param        limit query int false "Лимит для пагинации" default(20)
// @Success      200 {object}  ListCarTagsResponse  "Список тегов"
// @Router       /v1/cars/car-tags [get]
func (h *CarTagHandler) GetCarTags(c *gin.Context) {
	offset := int64(0)
	limit := int64(20)

	if o, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		offset = o
	}
	if l, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		limit = l
	}

	total, tags, err := h.getCarTags.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ListCarTagsResponse{
		Total: total,
		Data:  tags,
	})
}

// UpdateCarTag godoc
// @Summary      Обновление тега
// @Description  Обновляет информацию о теге по указанному ID
// @Tags         Car Tags
// @Accept       json
// @Produce      json
// @Param        id path int true "ID тега"
// @Param        request body UpdateCarTagRequest true "Новые данные для тега"
// @Success      200 {object}  CarTagResponse  "Тег успешно обновлен"
// @Security     BearerAuth
// @Router       /v1/cars/car-tags/{id} [put]
func (h *CarTagHandler) UpdateCarTag(c *gin.Context) {
	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID тега"))
		return
	}

	var req UpdateCarTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	tag, err := h.updateCarTag.Execute(c.Request.Context(), tagID, req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteCarTag godoc
// @Summary      Удаление тега
// @Description  Удаляет тег по указанному ID
// @Tags         Car Tags
// @Accept       json
// @Produce      json
// @Param        id path int true "ID тега"
// @Success      200 {object}  MessageResponse  "Тег успешно удален"
// @Security     BearerAuth
// @Router       /v1/cars/car-tags/{id} [delete]
func (h *CarTagHandler) DeleteCarTag(c *gin.Context) {
	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID тега"))
		return
	}

	err = h.deleteCarTag.Execute(c.Request.Context(), tagID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Тег успешно удален"})
}
