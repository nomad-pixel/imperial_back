package car

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/car_usecases"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CarHandler struct {
	createCar  usecasePorts.CreateCarUsecase
	deleteCar  usecasePorts.DeleteCarUsecase
	updateCar  usecasePorts.UpdateCarUsecase
	getCarById usecasePorts.GetCarByIdUsecase
	getCars    usecasePorts.GetListCarsUsecase
}

func NewCarHandler(
	createCar usecasePorts.CreateCarUsecase,
	deleteCar usecasePorts.DeleteCarUsecase,
	updateCar usecasePorts.UpdateCarUsecase,
	getCarById usecasePorts.GetCarByIdUsecase,
	getCars usecasePorts.GetListCarsUsecase,
) *CarHandler {
	return &CarHandler{
		createCar:  createCar,
		deleteCar:  deleteCar,
		updateCar:  updateCar,
		getCarById: getCarById,
		getCars:    getCars,
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

	c.JSON(http.StatusCreated, createdCar)
}

// DeleteCar godoc
// @Summary      Удаление автомобиля
// @Description  Удаляет автомобиль по указанному ID
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        id path int true "ID автомобиля для удаления"
// @Success      200 {object}  MessageResponse  "Автомобиль успешно удален"
// @Security     BearerAuth
// @Router       /v1/cars/{id} [delete]
func (h *CarHandler) DeleteCar(c *gin.Context) {
	carId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID автомобиля"))
		return
	}

	err = h.deleteCar.Execute(c.Request.Context(), int64(carId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := MessageResponse{
		Message: "Автомобиль успешно удален",
	}
	c.JSON(http.StatusOK, response)
}

// UpdateCar godoc
// @Summary      Обновление автомобиля
// @Description  Обновляет информацию об автомобиле по указанному ID
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        id path int true "ID автомобиля для обновления"
// @Param        request body UpdateCarRequest true "Данные для обновления автомобиля"
// @Success      200 {object}  CarResponse  "Автомобиль успешно обновлен"
// @Security     BearerAuth
// @Router       /v1/cars/{id} [put]
func (h *CarHandler) UpdateCar(c *gin.Context) {
	carId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID автомобиля"))
		return
	}

	var req UpdateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	car := &entities.Car{
		ID:             int64(carId),
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

	updatedCar, err := h.updateCar.Execute(c.Request.Context(), car)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedCar)
}

// GetCarById godoc
// @Summary      Получение автомобиля по ID
// @Description  Возвращает информацию об автомобиле по указанному ID
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        id path int true "ID автомобиля для получения"
// @Success      200 {object}  CarResponse  "Информация об автомобиле"
// @Router       /v1/cars/{id} [get]
func (h *CarHandler) GetCarByID(c *gin.Context) {
	carId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите ID автомобиля"))
		return
	}

	car, err := h.getCarById.Execute(c.Request.Context(), int64(carId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, car)
}

// GetCars godoc
// @Summary      Получение списка автомобилей
// @Description  Возвращает список автомобилей с возможностью фильтрации и пагинации
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        offset query int false "Смещение для пагинации" default(0)
// @Param        limit query int false "Лимит для пагинации" default(20)
// @Param        name query string false "Фильтр по названию автомобиля"
// @Param        mark_id query int false "Фильтр по ID марки автомобиля"
// @Param        category_id query int false "Фильтр по ID категории автомобиля"
// @Success      200 {object}  ListCarsResponse  "Список автомобилей"
// @Router       /v1/cars [get]
func (h *CarHandler) ListCars(c *gin.Context) {
	offset := int64(0)
	limit := int64(20)
	name := ""
	var markID int64 = 0
	var categoryID int64 = 0

	if o, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		offset = o
	}
	if l, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		limit = l
	}
	name = c.DefaultQuery("name", "")
	if m, err := strconv.ParseInt(c.DefaultQuery("mark_id", "0"), 10, 64); err == nil {
		markID = m
	}
	if cat, err := strconv.ParseInt(c.DefaultQuery("category_id", "0"), 10, 64); err == nil {
		categoryID = cat
	}

	total, cars, err := h.getCars.Execute(c.Request.Context(), offset, limit, name, markID, categoryID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ListCarsResponse{
		Total: total,
		Data:  cars,
	})
}
