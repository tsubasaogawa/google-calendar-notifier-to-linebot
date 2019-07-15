package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Cal struct{
	cred string
	maxResults int64
	srv *calendar.Service
	Plans []Plan
}

const (
	// format
	dtfmt = "2006-01-02T00:00:00.000Z"
)

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

	(*c).srv = srv
}

func (c *Cal) Retrieve(days int, onlyPubItem bool) {
	now := time.Now()
	st := now.Format(dtfmt)
	ed := now.AddDate(0, 0, days).Format(dtfmt)
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

			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}

			(*c).Plans = append((*c).Plans, Plan{
				date: date,
				title: item.Summary,
			})
			fmt.Printf("%+v\n\n", item)
		}
	}
}
