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

	log.Fatal(router.Run(":8085"))
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
		origins := os.Getenv("ALLOWED_ORIGINS")
		originsSlice := strings.Split(origins, ",")

		// Create regex patterns for each origin that allow for dynamic subdomains
		allowedOriginPatterns := make([]*regexp.Regexp, len(originsSlice))
		for i, origin := range originsSlice {
			// Escape special regex characters in the origin
			pattern := regexp.QuoteMeta(origin)

			// Allow an optional dynamic subdomain of the form:
			// {any-subdomain}--{number}.
			// This regex ensures that only the specified domains (and their subdomains) are allowed,
			// while permitting dynamic numbered variations.
			pattern = strings.Replace(pattern, "://", "://(?:[^.]+--\\d+\\.)?", 1)

			// Ensure full string match to prevent partial matches
			allowedOriginPatterns[i] = regexp.MustCompile("^" + pattern + "$")
		}

		fmt.Println(allowedOriginPatterns)

		return cors.New(cors.Config{
			AllowOriginFunc: func(origin string) bool {
				for _, pattern := range allowedOriginPatterns {
					if pattern.MatchString(origin) {
						return true
					}
				}
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
			},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	case "staging":
		return cors.Default()
	default:
		return nil
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
