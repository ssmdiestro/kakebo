package controllers

import (
	"kakebo/internal/dto"
	"kakebo/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRecord(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RecordRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := svc.NewRecord(c.Request.Context(), req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Record created successfully"})
	}
}

func GetRecordByDate(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")
		if _, err := time.Parse("2006-01-02", date); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Usa YYYY-MM-DD"})
			return
		}
		rec, err := svc.GetRecordByDate(c.Request.Context(), date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(rec) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No se encontró ningún registro para esa fecha"})
			return
		}
		c.JSON(http.StatusOK, rec)
	}
}
