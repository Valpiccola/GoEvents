package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	assert "github.com/go-playground/assert/v2"
)

func TestRecordEvent(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/record_event", RecordEvent)

	// Assert that function returns 400 when JSON is invalid
	t.Run(fmt.Sprintf("Test invalid JSON"), func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/record_event", bytes.NewBufferString(`{"Ip": 999}`))
		r.ServeHTTP(w, c.Request)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Assert that function returns 200 when passing a valid json but not country
	t.Run(fmt.Sprintf("Test valid JSON with deep == false so no country should be extracted"), func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/record_event", bytes.NewBufferString(`{"Deep": false}`))
		r.ServeHTTP(w, c.Request)
		// Select latest record from table events
		var country string
		_ = Db.QueryRow(`
			SELECT details#>>'{IpData,country}'
			FROM test.event
			ORDER BY created_at DESC
			LIMIT 1;`).Scan(&country)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, len(country), 0)
	})

	// Assert that function returns country when passing a valid json and deep is true
	t.Run(fmt.Sprintf("Test valid JSON with deep == true so country should be extracted"), func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/record_event", bytes.NewBufferString(`{"Deep": true}`))
		c.Request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
		r.ServeHTTP(w, c.Request)

		// Select latest record from table events
		var country string
		var browser string
		_ = Db.QueryRow(`
			SELECT
				details#>>'{IpData,country}',
				details#>>'{UserAgentData,Name}'
			FROM test.event
			ORDER BY created_at DESC
			LIMIT 1;`).Scan(&country, &browser)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, len(country), 2)
		assert.Equal(t, browser, "Chrome")
	})

	_, _ = Db.Exec(`DELETE FROM test.event WHERE 1 = 1;`)
}
