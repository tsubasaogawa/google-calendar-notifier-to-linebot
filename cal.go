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

func (c *Cal) Retrieve() {
	t := time.Now().Format(time.RFC3339)
	events, err := c.srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults((*c).maxResults).OrderBy("startTime").Do()
	if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
			fmt.Println("No upcoming events found.")
	} else {
			for _, item := range events.Items {
					date := item.Start.DateTime
					if date == "" {
							date = item.Start.Date
					}
					// fmt.Printf("%v (%v)\n", item.Summary, date)
					(*c).Plans = append((*c).Plans, Plan{
						date: date,
						title: item.Summary,
					})
			}
	}
}