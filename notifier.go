package main

import (
	"fmt"
)

const (
	// cred is a client secret file name.
	cred = "credentials.json"

	// maxResults is the maximum number of results.
	maxResults = 10

	// days
	days = 7

	onlyPubItem = true
)

func main() {
	cal := &Cal{}
	cal.NewCal(cred, maxResults)
	cal.Retrieve(days, onlyPubItem)

	fmt.Printf("%v", (*cal).Plans)
}