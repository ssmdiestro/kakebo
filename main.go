package main

import (
	"fmt"
	"kakebo/dto"
	"kakebo/repository"
	"kakebo/service"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI no definido")
	}

	if err := repository.SetMongoConnection(mongoURI); err != nil {
		log.Fatal("Error conectando a Mongo: ", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // <--- esta línea lo hace todo
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hola desde Gin!")
	})

	r.POST("/newRecord", func(c *gin.Context) {
		var recordRequest dto.RecordRequest
		if err := c.ShouldBindJSON(&recordRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := service.NewRecord(recordRequest); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Record created successfully"})
	})

	r.GET("/getRecord/:date", func(c *gin.Context) {
		date := c.Param("date")
		if _, err := time.Parse("2006-01-02", date); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Usa YYYY-MM-DD"})
			return
		}

		rec := service.GetRecordByDate(date)
		if len(rec) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No se encontró ningún registro para esa fecha"})
			return
		}
		c.JSON(http.StatusOK, rec)
	})

	r.GET("/getDayReport/:date", func(c *gin.Context) {
		date := c.Param("date")
		if _, err := time.Parse("2006-01-02", date); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Usa YYYY-MM-DD"})
			return
		}

		rec := service.GetDayReport(date)
		if rec.Total == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No se encontró ningún registro para esa fecha"})
			return
		}
		c.JSON(http.StatusOK, rec)
	})

	r.GET("/getWeek/:date", func(c *gin.Context) {
		date := c.Param("date")
		if _, err := time.Parse("2006-01-02", date); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Usa YYYY-MM-DD"})
			return
		}

		w, m, s, e, _ := service.WeekNumberInCustomMonth(date, time.Local)
		weekLimits := dto.WeekLimits{
			Week:      w,
			Month:     m,
			StartDate: s.Format("2006-01-02"),
			EndDate:   e.Format("2006-01-02"),
		}
		c.JSON(http.StatusOK, weekLimits)
	})

	r.GET("/getCategories", func(c *gin.Context) {
		categories := service.GetCategories()
		c.JSON(http.StatusOK, categories)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
