/**

######## Define the Event struct
This struct represents an event with fields such as Cookie, Referrer,
Page, Event_name, and other relevant information.
It also includes fields for storing IP-related information and parsed User-Agent data.

######## RecordEvent
This function handles incoming HTTP requests and does the following:
a. Bind the request JSON data to the Event struct.
b. Extract the client IP address and User-Agent header from the request.
c. If the Deep field is true, fetch IP-related details using GetIpDetails() function and parse User-Agent data.
d. Marshal the Event struct into JSON format.
e. Insert the JSON data into the event table in the database.
f. Return an appropriate HTTP response.

######## GetIpDetails
This function takes an IP address as input and returns IP-related information using the IPInfo package.
It requires an API token, which is set using the IPINFO_TOKEN environment variable.

*/

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

	q_event := `
		INSERT INTO %s.event (created_at, details)
		VALUES (current_timestamp, $1);
	`
	q_event_ready := fmt.Sprintf(q_event, os.Getenv("DB_SCHEMA"))
	_, err = Db.Exec(q_event_ready, b_event)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Error saving event in db",
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
