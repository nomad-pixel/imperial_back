package celebrity

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/celebrity"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type CelebrityHandler struct {
	createCelebrityUsecase      usecasePorts.CreateCelebrityUsecase
	celebrityUploadImageUsecase usecasePorts.UploadCelebrityImageUsecase
}

func NewCelebrityHandler(
	createCelebrityUsecase usecasePorts.CreateCelebrityUsecase,
	celebrityUploadImageUsecase usecasePorts.UploadCelebrityImageUsecase,
) *CelebrityHandler {
	return &CelebrityHandler{
		createCelebrityUsecase:      createCelebrityUsecase,
		celebrityUploadImageUsecase: celebrityUploadImageUsecase,
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
