package celebrity

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/celebrity"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CelebrityHandler struct {
	createCelebrityUsecase      usecasePorts.CreateCelebrityUsecase
	celebrityUploadImageUsecase usecasePorts.UploadCelebrityImageUsecase
	getCelebrityByIdUsecase     usecasePorts.GetCelebrityByIdUsecase
	listCelebritiesUsecase      usecasePorts.ListCelebritiesUsecase
	updateCelebrityUsecase      usecasePorts.UpdateCelebrityUsecase
	deleteCelebrityUsecase      usecasePorts.DeleteCelebrityUsecase
}

func NewCelebrityHandler(
	createCelebrityUsecase usecasePorts.CreateCelebrityUsecase,
	celebrityUploadImageUsecase usecasePorts.UploadCelebrityImageUsecase,
	getCelebrityByIdUsecase usecasePorts.GetCelebrityByIdUsecase,
	listCelebritiesUsecase usecasePorts.ListCelebritiesUsecase,
	updateCelebrityUsecase usecasePorts.UpdateCelebrityUsecase,
	deleteCelebrityUsecase usecasePorts.DeleteCelebrityUsecase,
) *CelebrityHandler {
	return &CelebrityHandler{
		createCelebrityUsecase:      createCelebrityUsecase,
		celebrityUploadImageUsecase: celebrityUploadImageUsecase,
		getCelebrityByIdUsecase:     getCelebrityByIdUsecase,
		listCelebritiesUsecase:      listCelebritiesUsecase,
		updateCelebrityUsecase:      updateCelebrityUsecase,
		deleteCelebrityUsecase:      deleteCelebrityUsecase,
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

	celebrity, err := h.createCelebrityUsecase.Execute(c.Request.Context(), req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(201, celebrity)
}

// UploadCelebrityImage godoc
// @Summary Upload an image for a celebrity
// @Description Upload an image file for the specified celebrity
// @Tags Celebrities
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Celebrity ID"
// @Param image formData file true "Image file"
// @Success 200 {object} entities.Celebrity
// @Router /v1/celebrities/{id}/image [put]
// @Security     BearerAuth
func (h *CelebrityHandler) UploadCelebrityImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid celebrity ID"))
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Image file is required"))
		return
	}

	fileData, err := file.Open()
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeInternal, "Failed to open image file"))
		return
	}
	defer fileData.Close()

	imageBytes, err := io.ReadAll(fileData)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeInternal, "Failed to read image file"))
		return
	}

	celebrity, err := h.celebrityUploadImageUsecase.Execute(c.Request.Context(), id, imageBytes, file.Filename)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, celebrity)
}

// GetCelebrityByID godoc
// @Summary Get celebrity by ID
// @Description Get detailed information about a celebrity by ID
// @Tags Celebrities
// @Accept json
// @Produce json
// @Param id path int true "Celebrity ID"
// @Success 200 {object} entities.Celebrity
// @Router /v1/celebrities/{id} [get]
// @Security     BearerAuth
func (h *CelebrityHandler) GetCelebrityByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid celebrity ID"))
		return
	}

	celebrity, err := h.getCelebrityByIdUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, celebrity)
}

// ListCelebrities godoc
// @Summary List celebrities
// @Description Get a paginated list of celebrities
// @Tags Celebrities
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(20)
// @Success 200 {object} ListCelebritiesResponse
// @Router /v1/celebrities [get]
// @Security     BearerAuth
func (h *CelebrityHandler) ListCelebrities(c *gin.Context) {
	offset := int64(0)
	limit := int64(20)

	if o, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		offset = o
	}
	if l, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		limit = l
	}

	total, celebrities, err := h.listCelebritiesUsecase.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, ListCelebritiesResponse{
		Total: total,
		Data:  celebrities,
	})
}

// UpdateCelebrity godoc
// @Summary Update celebrity
// @Description Update celebrity information by ID
// @Tags Celebrities
// @Accept json
// @Produce json
// @Param id path int true "Celebrity ID"
// @Param celebrity body UpdateCelebrityRequest true "Celebrity data"
// @Success 200 {object} entities.Celebrity
// @Router /v1/celebrities/{id} [put]
// @Security     BearerAuth
func (h *CelebrityHandler) UpdateCelebrity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid celebrity ID"))
		return
	}

	var req UpdateCelebrityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	updatedCelebrity, err := h.updateCelebrityUsecase.Execute(c.Request.Context(), id, req.Name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, updatedCelebrity)
}

// DeleteCelebrity godoc
// @Summary Delete celebrity
// @Description Delete a celebrity by ID
// @Tags Celebrities
// @Accept json
// @Produce json
// @Param id path int true "Celebrity ID"
// @Success 200 {object} MessageResponse
// @Router /v1/celebrities/{id} [delete]
// @Security     BearerAuth
func (h *CelebrityHandler) DeleteCelebrity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid celebrity ID"))
		return
	}

	err = h.deleteCelebrityUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, MessageResponse{
		Message: "Celebrity successfully deleted",
	})
}
