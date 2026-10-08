package postinteraction

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler{
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateInteractionPost(c *gin.Context) {
	var input InterationPostInput

	userID, _ := c.Get("user_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"jsonReq": err.Error(),
		})
		return
	}

	post, err := h.service.CreateInteractionPost(input, uint(userID.(float64)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}