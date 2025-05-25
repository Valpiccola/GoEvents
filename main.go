package main

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	DbUrl string
	Db    sql.DB
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	Db = *SetUpDb()
	defer Db.Close()

	router := gin.New()

	corsConfig := getCORSConfig()
	if corsConfig != nil {
		router.Use(corsConfig)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
		log.Info("Server is running")
	})

	router.POST("/record_event", RecordEvent)
	router.GET("/health", healthCheckHandler)

	log.Fatal(router.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func SetUpDb() (db *sql.DB) {

	DbUrl = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?client_encoding=utf8",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", DbUrl)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Unsucessfully connected with db",
		}).Error(err)
		panic(err)
	}

	if err = db.Ping(); err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Unsucessfully connected with db",
		}).Error(err)
		panic(err)
	}

	fmt.Println("Successfully connected to db")
	return
}

func getCORSConfig() gin.HandlerFunc {
	env := os.Getenv("ENV")
	switch env {
	case "production":
		allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
		allowedPatterns := strings.Split(os.Getenv("ALLOWED_PATTERNS"), ",")

		return cors.New(cors.Config{
			AllowOriginFunc: func(origin string) bool {
				fmt.Printf("CORS: Checking origin: '%s'\n", origin)

				// Check exact matches
				for _, allowedOrigin := range allowedOrigins {
					allowedOrigin = strings.TrimSpace(allowedOrigin)
					if allowedOrigin == origin {
						fmt.Printf("CORS: Exact match found for: '%s'\n", allowedOrigin)
						return true
					}
				}

				// Check regex patterns - fix the pattern matching
				for _, pattern := range allowedPatterns {
					pattern = strings.TrimSpace(pattern)
					if pattern != "" {

						// The pattern is already properly escaped, just use it directly
						if matched, err := regexp.MatchString("^"+pattern+"$", origin); err == nil && matched {
							return true
						} else if err != nil {
							fmt.Printf("CORS: Regex error for pattern '%s': %v\n", pattern, err)
						} else {
							fmt.Printf("CORS: Pattern '%s' did not match origin '%s'\n", pattern, origin)
						}
					}
				}

				fmt.Printf("CORS: No match found for origin: '%s'\n", origin)
				return false
			},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
			AllowHeaders: []string{
				"Accept",
				"Accept-Language",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				"Cache-Control",
				"X-Requested-With",
				"Origin",
				"sentry-trace",
				"baggage",
			},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	case "staging":
		fmt.Println("Environment: staging - using default CORS")
		return cors.Default()
	default:
		fmt.Println("Environment: development - allowing all origins")
		return cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
			AllowHeaders: []string{
				"Accept",
				"Accept-Language",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				"Cache-Control",
				"X-Requested-With",
				"Origin",
				"sentry-trace",
				"baggage",
			},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}
}

func healthCheckHandler(c *gin.Context) {
	if err := Db.Ping(); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"status": "error", "message": "Database is disconnected"},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"status": "success", "message": "API is healthy"},
	)
}
