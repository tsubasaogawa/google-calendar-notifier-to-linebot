package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/tsubasaogawa/linebot-publisher-layer-go"
)

const (
	// cred is a client secret file name.
	cred = "credentials.json"

	// maxResults is the maximum number of results.
	maxResults = 10

	// days is target days.
	days = 1

	// onlyPubItem is what obtains public items only if true.
	onlyPubItem = true
)

// Env is ...
type Env struct {
	Cred        string
	MaxResults  int64
	Days        int
	OnlyPubItem bool
	ToID        string
	AccessToken string
}

// Event is from lambda.
type Event struct {
	// :
}

// Response is returned by the function.
type Response struct {
	Num int `json:"Num"`
}

// notifier is
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

	// fmt.Printf("%v", (*cal).Plans)

	message := "今日の予定だよ\n"
	for _, plan := range plans {
		// fmt.Printf("Date: %s Title: %s\n", plan.date, plan.title)
		message += fmt.Sprintf("  %s %s\n", plan.date, plan.title)
	}
	linebot.Publish(env.ToID, message, false)

	return Response{Num: len(plans)}, nil
}

// getEnv is
func getEnv() Env {
	env := Env{}
	envconfig.Process("", &env)
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

	return env
}

func main() {
	lambda.Start(notifier)
}
