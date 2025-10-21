// backend/internal/http/router.go
package http

import (
	"net/http"
	"time"

	"kakebo/internal/http/controllers"
	"kakebo/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(svc *service.Service) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server Online!")
	})

	v1 := r.Group("/")

	// --- Categories ---
	v1.GET("/getCategories", controllers.GetCategories(svc))

	// --- RECORDS ---
	v1.POST("/newRecord", controllers.NewRecord(svc))
	v1.GET("/getRecord/:date", controllers.GetRecordByDate(svc))

	// --- REPORTS ---
	v1.GET("/getDayReport/:date", controllers.GetDayReport(svc))
	v1.GET("/getWeekReport", controllers.GetWeekReport(svc))
	v1.GET("/getWeekDays", controllers.GetWeekDays(svc))
	v1.GET("/getWeek/:date", controllers.GetWeekLimits(svc))

	return r
}
