package driver

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/driver"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type DriverHandler struct {
	createDriverUsecase      usecasePorts.CreateDriverUsecase
	getDriverByIdUsecase     usecasePorts.GetDriverByIdUsecase
	listDriversUsecase       usecasePorts.ListDriversUsecase
	updateDriverUsecase      usecasePorts.UpdateDriverUsecase
	deleteDriverUsecase      usecasePorts.DeleteDriverUsecase
	uploadDriverPhotoUsecase usecasePorts.UploadDriverPhotoUsecase
}

func NewDriverHandler(
	createDriverUsecase usecasePorts.CreateDriverUsecase,
	getDriverByIdUsecase usecasePorts.GetDriverByIdUsecase,
	listDriversUsecase usecasePorts.ListDriversUsecase,
	updateDriverUsecase usecasePorts.UpdateDriverUsecase,
	deleteDriverUsecase usecasePorts.DeleteDriverUsecase,
	uploadDriverPhotoUsecase usecasePorts.UploadDriverPhotoUsecase,
) *DriverHandler {
	return &DriverHandler{
		createDriverUsecase:      createDriverUsecase,
		getDriverByIdUsecase:     getDriverByIdUsecase,
		listDriversUsecase:       listDriversUsecase,
		updateDriverUsecase:      updateDriverUsecase,
		deleteDriverUsecase:      deleteDriverUsecase,
		uploadDriverPhotoUsecase: uploadDriverPhotoUsecase,
	}
}

func (h *DriverHandler) CreateDriver(c *gin.Context) {
	var req CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	driver, err := h.createDriverUsecase.Execute(c.Request.Context(), req.FullName, req.About, req.ExperienceYears)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(201, driver)
}

func (h *DriverHandler) GetDriverByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid driver ID"))
		return
	}

	driver, err := h.getDriverByIdUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, driver)
}

func (h *DriverHandler) ListDrivers(c *gin.Context) {
	offset := int64(0)
	limit := int64(20)

	if o, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		offset = o
	}
	if l, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		limit = l
	}

	total, drivers, err := h.listDriversUsecase.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, ListDriversResponse{
		Total: total,
		Data:  drivers,
	})
}

func (h *DriverHandler) UpdateDriver(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid driver ID"))
		return
	}

	var req UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	updatedDriver, err := h.updateDriverUsecase.Execute(c.Request.Context(), id, req.FullName, req.About, req.ExperienceYears)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, updatedDriver)
}

func (h *DriverHandler) DeleteDriver(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid driver ID"))
		return
	}

	err = h.deleteDriverUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Driver deleted successfully"})
}

func (h *DriverHandler) UploadDriverPhoto(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid driver ID"))
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Photo file is required"))
		return
	}

	fileData, err := file.Open()
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeInternal, "Failed to open photo file"))
		return
	}
	defer fileData.Close()

	imageBytes, err := io.ReadAll(fileData)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeInternal, "Failed to read photo file"))
		return
	}

	driver, err := h.uploadDriverPhotoUsecase.Execute(c.Request.Context(), id, imageBytes, file.Filename)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, driver)
}
