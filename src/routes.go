package main

import (
	"fmt"
)

func s(a string) *string {
	tmp := a
	return &tmp
}

// PongRoute ...
func PongRoute(event eventInfo) (bool, *string) {
	sent, _ := event.message.Timestamp.Parse()
	return false, s(fmt.Sprintf("%fs", sent.Sub(event.received).Seconds()))
}
