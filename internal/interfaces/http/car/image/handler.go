package car

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/car"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CarImageHandler struct {
	createCarImage   usecasePorts.CreateCarImageUsecase
	deleteCarImage   usecasePorts.DeleteCarImageUsecase
	getCarImagesList usecasePorts.GetCarImagesListUsecase
}

func NewCarImageHandler(
	createCarImage usecasePorts.CreateCarImageUsecase,
	deleteCarImage usecasePorts.DeleteCarImageUsecase,
	getCarImagesList usecasePorts.GetCarImagesListUsecase,
) *CarImageHandler {
	return &CarImageHandler{
		createCarImage:   createCarImage,
		deleteCarImage:   deleteCarImage,
		getCarImagesList: getCarImagesList,
	}
}

// CreateCarImage godoc
// @Summary      Загрузка изображения автомобиля
// @Description  Загружает изображение для указанного автомобиля
// @Tags         Car images
// @Accept       multipart/form-data
// @Produce      json
// @Param        car_id path int true "ID автомобиля"
// @Param        image formData file true "Изображение автомобиля"
// @Security     BearerAuth
// @Router       /v1/cars/{car_id}/images [post]
// @Success      201 {object}  CarImageResponse  "Марка успешно создана"
func (h *CarImageHandler) CreateCarImage(c *gin.Context) {
	carID, err := strconv.ParseInt(c.Param("car_id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите корректный ID автомобиля"))
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}

	image, err := h.createCarImage.Execute(c.Request.Context(), carID, fileData, fileHeader.Filename)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, image)
}

// DeleteCarImage godoc
// @Summary      Удаление изображения автомобиля
// @Description  Удаляет изображение автомобиля по указанному пути
// @Tags         Car images
// @Accept       json
// @Produce      json
// @Param        id path int true "ID изображения"
// @Security     BearerAuth
// @Router       /v1/cars/images/{id} [delete]
// @Success      200 {object}  MessageResponse  "Изображение успешно удалено"
func (h *CarImageHandler) DeleteCarImage(c *gin.Context) {
	imageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите корректный ID изображения"))
		return
	}

	err = h.deleteCarImage.Execute(c.Request.Context(), imageID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение успешно удалено"})
}

// GetCarImagesList godoc
// @Summary      Получение списка изображений автомобиля
// @Description  Возвращает список изображений для указанного автомобиля
// @Tags         Car images
// @Accept       json
// @Produce      json
// @Param        car_id path int true "ID автомобиля"
// @Param        offset query int false "Смещение для пагинации" default(0)
// @Param        limit query int false "Лимит для пагинации" default(20)
// @Security     BearerAuth
// @Router       /v1/cars/{car_id}/images [get]
// @Success      200 {object}  ListCarImagesResponse  "Список изображений успешно получен"
func (h *CarImageHandler) GetCarImagesList(c *gin.Context) {
	carID, err := strconv.ParseInt(c.Param("car_id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Укажите корректный ID автомобиля"))
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Неверный формат смещения"))
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		_ = c.Error(errors.New(errors.ErrCodeValidation, "Неверный формат лимита"))
		return
	}

	total, images, err := h.getCarImagesList.Execute(c.Request.Context(), carID, offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ListCarImagesResponse{
		Total: total,
		Data:  images,
	})

}
