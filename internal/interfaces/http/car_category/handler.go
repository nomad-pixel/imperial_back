package car_category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/car"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CarCategoryHandler struct {
	createCarCategory usecasePorts.CreateCarCategoryUsecase
	getCarCategory    usecasePorts.GetCarCategoryUsecase
	getCarCategories  usecasePorts.GetCarCategoriesListUsecase
	updateCarCategory usecasePorts.UpdateCarCategoryUsecase
	deleteCarCategory usecasePorts.DeleteCarCategoryUsecase
}

func NewCarCategoryHandler(
	createCarCategory usecasePorts.CreateCarCategoryUsecase,
	getCarCategory usecasePorts.GetCarCategoryUsecase,
	getCarCategories usecasePorts.GetCarCategoriesListUsecase,
	updateCarCategory usecasePorts.UpdateCarCategoryUsecase,
	deleteCarCategory usecasePorts.DeleteCarCategoryUsecase,
) *CarCategoryHandler {
	return &CarCategoryHandler{
		createCarCategory: createCarCategory,
		getCarCategory:    getCarCategory,
		getCarCategories:  getCarCategories,
		updateCarCategory: updateCarCategory,
		deleteCarCategory: deleteCarCategory,
	}
}

// CreateCarCategory godoc
// @Summary      Создание новой категории машины
// @Description  Создает новую категорию с указанным названием
// @Tags         Car Categories
// @Accept       json
// @Produce      json
// @Param        request body CreateCarCategoryRequest true "Данные для создания категории"
// @Success      201 {object}  CarCategoryResponse  "Категория успешно создана"
// @Security     BearerAuth
// @Router       /v1/cars/car-categories [post]
func (h *CarCategoryHandler) CreateCarCategory(c *gin.Context) {
	var req CreateCarCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	category, err := h.createCarCategory.Execute(c.Request.Context(), req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCarCategory godoc
// @Summary      Получение категории по ID
// @Description  Возвращает информацию о категории по указанному ID
// @Tags         Car Categories
// @Accept       json
// @Produce      json
// @Param        id path int true "ID категории"
// @Success      200 {object}  CarCategoryResponse  "Информация о категории"
// @Security     BearerAuth
// @Router       /v1/cars/car-categories/{id} [get]
func (h *CarCategoryHandler) GetCarCategory(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID категории"))
		return
	}

	category, err := h.getCarCategory.Execute(c.Request.Context(), categoryID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCarCategories godoc
// @Summary      Получение списка категорий
// @Description  Возвращает список всех категорий с поддержкой пагинации
// @Tags         Car Categories
// @Accept       json
// @Produce      json
// @Param        offset query int false "Смещение для пагинации" default(0)
// @Param        limit query int false "Лимит для пагинации" default(20)
// @Success      200 {object}  ListCarCategoriesResponse  "Список категорий"
// @Security     BearerAuth
// @Router       /v1/cars/car-categories [get]
func (h *CarCategoryHandler) GetCarCategories(c *gin.Context) {
	offset := int64(0)
	limit := int64(20)

	if o, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		offset = o
	}
	if l, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		limit = l
	}

	total, categories, err := h.getCarCategories.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ListCarCategoriesResponse{
		Total: total,
		Data:  categories,
	})
}

// UpdateCarCategory godoc
// @Summary      Обновление категории
// @Description  Обновляет информацию о категории по указанному ID
// @Tags         Car Categories
// @Accept       json
// @Produce      json
// @Param        id path int true "ID категории"
// @Param        request body UpdateCarCategoryRequest true "Новые данные для категории"
// @Success      200 {object}  CarCategoryResponse  "Категория успешно обновлена"
// @Security     BearerAuth
// @Router       /v1/cars/car-categories/{id} [put]
func (h *CarCategoryHandler) UpdateCarCategory(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID категории"))
		return
	}

	var req UpdateCarCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	category, err := h.updateCarCategory.Execute(c.Request.Context(), categoryID, req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCarCategory godoc
// @Summary      Удаление категории
// @Description  Удаляет категорию по указанному ID
// @Tags         Car Categories
// @Accept       json
// @Produce      json
// @Param        id path int true "ID категории"
// @Success      200 {object}  MessageResponse  "Категория успешно удалена"
// @Security     BearerAuth
// @Router       /v1/cars/car-categories/{id} [delete]
func (h *CarCategoryHandler) DeleteCarCategory(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID категории"))
		return
	}

	err = h.deleteCarCategory.Execute(c.Request.Context(), categoryID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Категория успешно удалена"})
}
