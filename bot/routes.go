package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"haze.sh/mairu/api"
	game "haze.sh/mairu/gutil"
	str "haze.sh/mairu/strutil"
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

// StatusRoute ...
func StatusRoute(event eventInfo) (bool, *string) {
	args := event.arguments
	name := strings.Join(event.arguments[1:], " ")
	dg := event.sesh
	switch strings.ToLower(args[1]) {
	case "g":
	case "game":
		dg.UpdateStatus(0, name)
	case "listening":
	case "l":
		game.UpdateStatusSpecial(dg, false, name, game.TypeListening)
	case "streaming":
	case "s":
		dg.UpdateStreamingStatus(0, strings.Join(args[3:], " "), args[2])
	case "watching":
	case "w":
		game.UpdateStatusSpecial(dg, false, name, game.TypeWatching)
	}
	return true, nil
}
