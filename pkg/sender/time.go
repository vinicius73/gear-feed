package sender

import "time"

func CalculeSendInterval(count int) time.Duration {
	switch {
	case count >= 20:
		// we cant sent more than 20 messages per minute
		return time.Second * 3
	case count >= 10:
		return time.Second * 2
	default:
		return time.Second
	}
}
