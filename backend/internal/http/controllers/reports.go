package controllers

import (
	"kakebo/internal/dto"
	"kakebo/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDayReport(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")
		if _, err := time.Parse("2006-01-02", date); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Usa YYYY-MM-DD"})
			return
		}
		rec, err := svc.GetDayReport(c.Request.Context(), date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rec.Total == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No se encontró ningún registro para esa fecha"})
			return
		}
		c.JSON(http.StatusOK, rec)
	}
}

func GetWeekLimits(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date")
		if _, err := time.Parse("2006-01-02", date); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Usa YYYY-MM-DD"})
			return
		}
		w, m, s, e, err := service.WeekNumberInCustomMonth(date, time.Local)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp := dto.WeekLimits{
			Week:      w,
			Month:     m,
			StartDate: s.Format("2006-01-02"),
			EndDate:   e.Format("2006-01-02"),
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetWeekDays(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		month, _ := strconv.Atoi(c.Query("month"))
		week, _ := strconv.Atoi(c.Query("week"))
		year, _ := strconv.Atoi(c.Query("year"))

		weekDays, err := service.WeekDaysInCustomMonth(year, month, week, time.Local)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, weekDays)
	}
}

func GetWeekReport(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		week, _ := strconv.Atoi(c.Query("week"))
		month, _ := strconv.Atoi(c.Query("month"))
		year, _ := strconv.Atoi(c.Query("year"))
		weekReport, err := svc.GetWeekReport(c.Request.Context(), week, month, year)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, weekReport)
	}
}
