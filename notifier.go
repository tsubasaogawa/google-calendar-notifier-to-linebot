package main

import (
	"fmt"
)

const (
	// cred is a client secret file name.
	cred = "credentials.json"

	// maxResults is the maximum number of results.
	maxResults = 10
)

func main() {
	cal := &Cal{}
	cal.NewCal(cred, maxResults)
	cal.Retrieve()

	fmt.Printf("%v", (*cal).Plans)
}