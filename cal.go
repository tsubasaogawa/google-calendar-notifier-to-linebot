package main

// Calendar object.

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Cal is a calendar object.
type Cal struct {
	// credentials json file
	cred string
	// maximum number of results
	maxResults int64
	// google calendar service object
	srv *calendar.Service
	// obtained plan list
	Plans []Plan
}

const (
	// time format
	dtfmt = "2006-01-02T00:00:00.000Z"

	// time zone; default value is JST (+9)
	convTime = "+09:00"
)

// NewCal is that creates a new calendar object.
func (c *Cal) NewCal(cred string, maxResults int64) {
	(*c).cred = cred
	(*c).maxResults = maxResults

	c.createGoogleCalendar()
}

func (c *Cal) createGoogleCalendar() {
	b, err := ioutil.ReadFile((*c).cred)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := Client{}.Create(config)
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// Set service object
	(*c).srv = srv
}

// Retrieve gets plans from the calendar.
func (c *Cal) Retrieve(days int, onlyPubItem bool) {
	// Set target date range (today 0 AM - tomorrow 0 AM)
	now := time.Now()
	st := now.Format(dtfmt)
	ed := now.AddDate(0, 0, days).Format(dtfmt)
	// Get plans
	events, err := c.srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(st).TimeMax(ed).MaxResults((*c).maxResults).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next %d of the user's events: %v", (*c).maxResults, err)
	}

	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			if onlyPubItem && item.Visibility == "private" {
				continue
			}

			date := item.Start.Date
			if item.Start.DateTime != "" {
				date = c.formatDateTime(item.Start.DateTime)
			}

			(*c).Plans = append((*c).Plans, Plan{
				date:  date,
				title: item.Summary,
			})
			// fmt.Printf("%+v\n\n", item)
		}
	}
}

func (c *Cal) formatDateTime(date string) string {
	dt, _ := time.Parse("2006-01-02T15:04:05"+convTime, date)
	return dt.Format("2006-01-02 15:04")
}
