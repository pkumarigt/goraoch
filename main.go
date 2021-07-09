package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"goroach/quote"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
)

var (
	db     *gorm.DB
	router *gin.Engine
	log    zerolog.Logger
)

func main() {
	//sleep for 15 seconds. Quick fix for a situation when DB is slow to start
	time.Sleep(15 * time.Second)
	// Initialize Dependencies
	// Service Port, Database, Logger, Cache, Message Queue etc.
	router := gin.Default()
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
	// Database
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_SERVER"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_DATABASE"), os.Getenv("DB_PASSWORD"),
	)
	log.Info().Msg(dsn)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer db.Close()

	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("LOG", log)
	})

	// Boostrap services
	quoteSvc := &quote.QuoteService{}
	quoteSvc.Bootstrap(router)

	// --- Development Only ---
	setupQuotes(db)

	// Start the service
	router.GET("/health", healthsvc)
	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

func healthsvc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": health()})
}

// Utility function to populate some records into DB
// Should not be used in production.
// Production should use sql scripts to create DB with default tables and data
func setupQuotes(db *gorm.DB) {
	// check if table exists
	// if table exists, return
	if !db.HasTable(&quote.QuoteModel{}) {
		db.AutoMigrate(&quote.QuoteModel{})

		quotes := []quote.QuoteModel{
			{Author: "Gandhi", Quote: "The best way to find yourself is to lose yourself in the service of others."},
			{Author: "Duke Ellington", Quote: "A problem is a chance for you to do your best."},
			{Author: "Steve Prefontaine", Quote: "To give anything less than your best, is to sacrifice the gift."},
			{Author: "Peter Drucker", Quote: "The best way to predict the future is to create it."},
		}

		for i := range quotes {
			db.Create(&quotes[i])
		}
	}
}
