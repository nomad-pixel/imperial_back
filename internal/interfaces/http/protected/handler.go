package protected

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProtectedHandler struct {
}

func NewProtectedHandler() *ProtectedHandler {
	return &ProtectedHandler{}
}

func (h *ProtectedHandler) Me(c *gin.Context) {
	uid, _ := c.Get("user_id")
	c.JSON(http.StatusOK, gin.H{
		"message": "protected endpoint",
		"user_id": uid,
	})
}
