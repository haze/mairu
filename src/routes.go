package main

import (
	"fmt"
	"time"
	"math"
)

func s(a string) *string {
	tmp := a
	return &tmp
}

// PongRoute ...
func PongRoute(event eventInfo) (bool, *string) {
	sent, _ := event.message.Timestamp.Parse()
	elapsed := math.Abs(float64(time.Since(sent) / time.Millisecond))
	return false, s(fmt.Sprintf("%vms", elapsed))
}
