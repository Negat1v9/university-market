package paymentservice

import "time"

// checks that more time has passed since time t than minInterval
func checkTimeCreatedPayment(minInterval time.Duration, t time.Time) bool {
	return time.Since(t) >= minInterval
}
