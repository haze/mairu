package main

import (
	"fmt" // exports the S function
	"math"
	"strings"
	"time"

	"haze.sh/mairu/api"
	"haze.sh/mairu/util"
)

// WolframRoute ...
func WolframRoute(event eventInfo) (bool, *string) {
	before := strings.Join(event.arguments, " ")
	res, _ := wolfram.Ask(event.config.WolframAlphaKey, before)
	if *res == "Wolfram|Alpha did not understand your input" {
		return false, str.S(fmt.Sprintf("%s = ???", before))
	}
	return false, str.S(fmt.Sprintf("%s = %s", before, *res))
}

// PongRoute ...
func PongRoute(event eventInfo) (bool, *string) {
	sent, _ := event.message.Timestamp.Parse()
	elapsed := math.Abs(float64(time.Since(sent) / time.Millisecond))
	return false, str.S(fmt.Sprintf("%vms", elapsed))
}
