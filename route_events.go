package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/mileusna/useragent"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	Cookie        string                 `json:"Cookie"`
	Referrer      string                 `json:"Referrer"`
	Page          string                 `json:"Page"`
	Event_name    string                 `json:"Event_name"`
	UserID        string                 `json:"UserID"`
	Size          string                 `json:"Size"`
	Language      string                 `json:"Language"`
	Deep          bool                   `json:"Deep"`
	Details       map[string]interface{} `json:"Details"`
	Ip            string
	UserAgent     string
	IpData        *ipinfo.Core
	UserAgentData useragent.UserAgent
}

func RecordEvent(c *gin.Context) {

	var event Event

	if err := c.BindJSON(&event); err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed binding event to JSON",
		}).Error(err)
		c.String(http.StatusBadRequest, "KO")
		return
	}

	fmt.Println(event)

	event.Ip = c.ClientIP()
	event.UserAgent = c.Request.Header.Get("User-Agent")
	if event.Deep {
		event.IpData = GetIpDetails(event.Ip)
		event.UserAgentData = useragent.Parse(event.UserAgent)
	}

	b_event, err := json.Marshal(event)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed marshaling event in RecordEvent",
		}).Error(err)
		c.String(http.StatusBadRequest, "KO")
		return
	}

	q := fmt.Sprintf(`
		INSERT INTO %s.event (created_at, details)
		VALUES (current_timestamp, $1);
	`, os.Getenv("DB_SCHEMA"))

	stmt, err := Db.Prepare(q)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Error preparing query",
		}).Error(err)
		c.String(http.StatusBadRequest, "KO")
		return
	}

	_, err = stmt.Exec(b_event)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Error saving event in db",
		}).Error(err)
		c.String(http.StatusBadRequest, "KO")
		return
	}

	c.String(http.StatusOK, "OK")
}

func GetIpDetails(ip_address string) (info *ipinfo.Core) {
	client := ipinfo.NewClient(nil, nil, os.Getenv("IPINFO_TOKEN"))
	info, err := client.GetIPInfo(net.ParseIP(ip_address))
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Error parsing IP address",
		}).Error(err)
	}
	return
}
