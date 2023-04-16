package main

import (
	"os"

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
	router.Use(SetUpCORS())
	router.POST("/record_event", RecordEvent)

	router.Run(":8080")
}

func SetUpDb() (db *sql.DB) {

	DbUrl = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
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

func SetUpCORS() gin.HandlerFunc {

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		fmt.Println("THIS IS THE ORIGIN", c.Request.Header.Get("Origin"))

		if os.Getenv("ENV") == "production" {
			allowedOrigins := map[string]bool{
				"https://valpiccola.com":     true,
				"https://www.valpiccola.com": true,
			}

			if _, ok := allowedOrigins[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else if os.Getenv("ENV") == "staging" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

/* func SetUpCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		env := os.Getenv("ENV")

		if env == "production" {
			allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
			allowedOriginsList := strings.Split(allowedOriginsStr, ",")

			allowedOrigins := make(map[string]bool)
			for _, o := range allowedOriginsList {
				allowedOrigins[strings.TrimSpace(o)] = true
				fmt.Println(strings.TrimSpace(o))
			}

			fmt.Println(allowedOrigins)

			if _, ok := allowedOrigins[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin",
					origin)
			}
		} else if env == "staging" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, "+
				"Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"POST, OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} */
