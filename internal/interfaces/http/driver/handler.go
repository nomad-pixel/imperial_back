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

// CreateDriver godoc
// @Summary Create a new driver
// @Description Create a new driver with the provided details
// @Tags Drivers
// @Accept json
// @Produce json
// @Param driver body CreateDriverRequest true "Driver data"
// @Success 201 {object} entities.Driver
// @Router /v1/drivers [post]
// @Security     BearerAuth
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

// GetDriverByID godoc
// @Summary Get driver by ID
// @Description Get detailed information about a driver by ID
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path int true "Driver ID"
// @Success 200 {object} entities.Driver
// @Router /v1/drivers/{id} [get]
// @Security     BearerAuth
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

// ListDrivers godoc
// @Summary List drivers
// @Description Get a paginated list of drivers
// @Tags Drivers
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(20)
// @Success 200 {object} ListDriversResponse
// @Router /v1/drivers [get]
// @Security     BearerAuth
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

// UpdateDriver godoc
// @Summary Update driver
// @Description Update driver information by ID
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path int true "Driver ID"
// @Param driver body UpdateDriverRequest true "Driver data"
// @Success 200 {object} entities.Driver
// @Router /v1/drivers/{id} [put]
// @Security     BearerAuth
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

// DeleteDriver godoc
// @Summary Delete driver
// @Description Delete a driver by ID
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path int true "Driver ID"
// @Success 200 {object} map[string]string
// @Router /v1/drivers/{id} [delete]
// @Security     BearerAuth
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

// UploadDriverPhoto godoc
// @Summary Upload a photo for a driver
// @Description Upload a photo file for the specified driver
// @Tags Drivers
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Driver ID"
// @Param photo formData file true "Photo file"
// @Success 200 {object} entities.Driver
// @Router /v1/drivers/{id}/photo [put]
// @Security     BearerAuth
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
