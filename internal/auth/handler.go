package auth

import (
	"net/http"

	"github.com/LuisFernando-hub/go-social-api/internal/config"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	cfg *config.Config
}

func NewHandler(service *Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		cfg: cfg,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginAuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"jsonReq": err.Error(),
		})
		return
	}

	token, err := h.service.Login(h.cfg, input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, token)
}