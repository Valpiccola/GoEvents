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
		fmt.Println("Allowed Origins:", allowedOrigins)
		fmt.Println("Allowed Patterns:", allowedPatterns)
		return cors.New(cors.Config{
			AllowOriginFunc: func(origin string) bool {
				fmt.Println("Checking Origin:", origin)

				// Check exact matches
				for _, allowedOrigin := range allowedOrigins {
					allowedOrigin = strings.TrimSpace(allowedOrigin)
					if allowedOrigin == origin {
						fmt.Println("Exact match found:", allowedOrigin)
						return true
					}
				}

				// Check patterns
				for _, pattern := range allowedPatterns {
					pattern = strings.TrimSpace(pattern)
					if pattern != "" {
						// Escape dots in the pattern for proper regex matching
						escapedPattern := strings.ReplaceAll(pattern, ".", "\\.")
						if matched, err := regexp.MatchString("^"+escapedPattern+"$", origin); err == nil && matched {
							fmt.Println("Pattern match found:", pattern, "for origin:", origin)
							return true
						}
					}
				}

				fmt.Println("No match found for origin:", origin)
				return false
			},
			AllowMethods: []string{"POST", "OPTIONS", "GET"},
			AllowHeaders: []string{
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				"accept",
				"origin",
				"Cache-Control",
				"X-Requested-With",
				"sentry-trace",
				"baggage",
			},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	case "staging":
		return cors.Default()
	default:
		return cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"POST", "OPTIONS", "GET"},
			AllowHeaders: []string{
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				"accept",
				"origin",
				"Cache-Control",
				"X-Requested-With",
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
