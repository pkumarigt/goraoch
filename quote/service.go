package quote

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type QuoteService struct{}

// All the services should be protected by auth token
func (qs *QuoteService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/quote")
	qrg.GET("/", qs.GetAll)
	qrg.GET("/:id", qs.Get)
	qrg.POST("/", qs.New)
	qrg.PUT("/:id", qs.Update)
	qrg.DELETE("/:id", qs.Delete)
}

// Get all quotes
// Not efficient and shoule limit results to default size
func (qs *QuoteService) GetAll(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var quotes []QuoteModel
	db.Limit(os.Getenv("DEFAULT_PAGE_SIZE")).Find(&quotes)
	c.JSON(http.StatusOK, quotes)
}

// Get a quote based on the ID
func (qs *QuoteService) Get(c *gin.Context) {
	var quote QuoteModel
	db := c.MustGet("DB").(*gorm.DB)
	err := db.Where("id = ?", c.Param("id")).First(&quote).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// Add a new quote to the database
func (qs *QuoteService) New(c *gin.Context) {
	var qr QuoteRequest
	err := c.ShouldBindJSON(&qr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quote := QuoteModel{Author: qr.Author, Quote: qr.Quote}
	db := c.MustGet("DB").(*gorm.DB)
	db.Create(&quote)

	c.JSON(http.StatusOK, gin.H{"data": quote})
}

// Update a quote based on the ID
func (qs *QuoteService) Update(c *gin.Context) {
	var quote QuoteModel
	db := c.MustGet("DB").(*gorm.DB)
	err := db.Where("id = ?", c.Param("id")).First(&quote).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var qr QuoteRequest
	err2 := c.ShouldBindJSON(&qr)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	db.Model(&quote).Updates(qr)

	c.JSON(http.StatusOK, gin.H{"data": quote})
}

// Delete a quote based on the ID
func (qs *QuoteService) Delete(c *gin.Context) {
	var quote QuoteModel
	db := c.MustGet("DB").(*gorm.DB)
	err := db.Where("id = ?", c.Param("id")).First(&quote).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Delete(&quote)

	c.JSON(http.StatusOK, gin.H{"data": quote})
}
