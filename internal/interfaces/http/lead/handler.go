package lead

import (
	"strconv"

	"github.com/gin-gonic/gin"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases/lead"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type LeadHandler struct {
	createLeadUsecase  usecasePorts.CreateLeadUsecase
	getLeadByIdUsecase usecasePorts.GetLeadByIdUsecase
	listLeadsUsecase   usecasePorts.ListLeadsUsecase
	deleteLeadUsecase  usecasePorts.DeleteLeadUsecase
}

func NewLeadHandler(
	createLeadUsecase usecasePorts.CreateLeadUsecase,
	getLeadByIdUsecase usecasePorts.GetLeadByIdUsecase,
	listLeadsUsecase usecasePorts.ListLeadsUsecase,
	deleteLeadUsecase usecasePorts.DeleteLeadUsecase,
) *LeadHandler {
	return &LeadHandler{
		createLeadUsecase:  createLeadUsecase,
		getLeadByIdUsecase: getLeadByIdUsecase,
		listLeadsUsecase:   listLeadsUsecase,
		deleteLeadUsecase:  deleteLeadUsecase,
	}
}

// CreateLead godoc
// @Summary Create a new lead
// @Description Create a new lead with the provided details
// @Tags Leads
// @Accept json
// @Produce json
// @Param lead body CreateLeadRequest true "Lead data"
// @Success 201 {object} entities.Lead
// @Router /v1/leads [post]
// @Security     BearerAuth
func (h *LeadHandler) CreateLead(c *gin.Context) {
	var req CreateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат данных"))
		return
	}

	lead, err := h.createLeadUsecase.Execute(c.Request.Context(), req.FullName, req.Phone, req.StartDate, req.EndDate)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(201, lead)
}

// GetLeadByID godoc
// @Summary Get lead by ID
// @Description Get detailed information about a lead by ID
// @Tags Leads
// @Accept json
// @Produce json
// @Param id path int true "Lead ID"
// @Success 200 {object} entities.Lead
// @Router /v1/leads/{id} [get]
// @Security     BearerAuth
func (h *LeadHandler) GetLeadByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid lead ID"))
		return
	}

	lead, err := h.getLeadByIdUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, lead)
}

// ListLeads godoc
// @Summary List leads
// @Description Get a paginated list of leads
// @Tags Leads
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(20)
// @Success 200 {object} ListLeadsResponse
// @Router /v1/leads [get]
// @Security     BearerAuth
func (h *LeadHandler) ListLeads(c *gin.Context) {
	offset := int64(0)
	limit := int64(20)

	if o, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		offset = o
	}
	if l, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		limit = l
	}

	total, leads, err := h.listLeadsUsecase.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, ListLeadsResponse{
		Total: total,
		Data:  leads,
	})
}

// DeleteLead godoc
// @Summary Delete lead
// @Description Delete a lead by ID
// @Tags Leads
// @Accept json
// @Produce json
// @Param id path int true "Lead ID"
// @Success 200 {object} map[string]string
// @Router /v1/leads/{id} [delete]
// @Security     BearerAuth
func (h *LeadHandler) DeleteLead(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Invalid lead ID"))
		return
	}

	err = h.deleteLeadUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Lead deleted successfully"})
}
