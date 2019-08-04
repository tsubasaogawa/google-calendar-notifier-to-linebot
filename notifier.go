package main

// Notifier main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/tsubasaogawa/linebot-publisher-layer-go"
)

// Define default values if environment variable is not set.
const (
	// client secret file name
	cred = "credentials.json"

	// the maximum number of results
	maxResults = 10

	// target days
	days = 1

	// If true, notifier shows only public items.
	onlyPubItem = true

	// Custom text showed above obtained plans
	description = "今日の予定だよ"
)

// Env has properties of environment variable; used by envconfig
type Env struct {
	Cred        string
	MaxResults  int64
	Days        int
	OnlyPubItem bool
	Description string
	// LINE ID
	ToID string
	// LINE access token
	AccessToken string
}

// Event is lambda events.
type Event struct {
	// :
}

// Response is returned value by the function.
type Response struct {
	Num int `json:"Num"`
}

func notifier(event Event) (Response, error) {
	env := getEnv()

	// create calendar
	cal := &Cal{}
	cal.NewCal(env.Cred, env.MaxResults)

	cal.Retrieve(env.Days, env.OnlyPubItem)
	plans := (*cal).Plans
	if len(plans) == 0 {
		return Response{Num: 0}, nil
	}

	fmt.Printf("%v", (*cal).Plans)

	message := env.Description
	for _, plan := range plans {
		message += fmt.Sprintf("\n  %s %s", plan.date, plan.title)
	}
	linebot.Publish(env.ToID, env.AccessToken, message, false)

	return Response{Num: len(plans)}, nil
}

func getEnv() Env {
	env := Env{}
	envconfig.Process("", &env)

	// Set default value if env value is not set
	if env.Cred == "" {
		env.Cred = cred
	}
	if env.MaxResults == 0 {
		env.MaxResults = maxResults
	}
	if env.Days == 0 {
		env.Days = days
	}
	if env.OnlyPubItem == false {
		env.OnlyPubItem = onlyPubItem
	}
	if env.Description == "" {
		env.Description = description
	}

	return env
}

func main() {
	lambda.Start(notifier)
}
