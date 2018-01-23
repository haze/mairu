package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"haze.sh/mairu/api"
	"haze.sh/mairu/util"
)

// WolframRoute ...
func WolframRoute(event eventInfo) (bool, *string) {
	before := strings.Join(event.arguments[1:], " ")
	if event.arguments[0] == "?" {
		res, err := wolfram.AskSimple(event.config.WolframAlphaKey, before)
		if *res == "Wolfram|Alpha did not understand your input" {
			return false, str.S(fmt.Sprintf("%s = ??? :(", before))
		}
		if err != nil {
			return true, nil
		}
		return false, str.S(fmt.Sprintf("%s = %s", before, *res))
	}
	res, err := wolfram.AskAdvanced(event.config.WolframAlphaKey, before)
	if err != nil || res.Success == false {
		return true, nil
	}
	rfmt := "```Input: %s\nOutput: %s\nTook: %fms```"
	fmt.Printf("%+v\n", res)
	return false, str.S(fmt.Sprintf(rfmt, *res.GetInterpretation(), *res.GetResult(), res.Timing))
}

// PongRoute ...
func PongRoute(event eventInfo) (bool, *string) {
	sent, _ := event.message.Timestamp.Parse()
	elapsed := math.Abs(float64(time.Since(sent) / time.Millisecond))
	return false, str.S(fmt.Sprintf("%vms", elapsed))
}
