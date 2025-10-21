package controllers

import (
	"kakebo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategories(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats, err := svc.GetCategories(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cats)
	}
}
