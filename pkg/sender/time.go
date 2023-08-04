package sender

import "time"

func CalculeSendInterval(count int) time.Duration {
	// we cant sent more than 20 messages per minute
	if count >= 20 {
		return time.Second * 3
	}

	dur := time.Duration(30/count) * time.Second

	if dur > time.Second*2 {
		return time.Second * 1
	}

	return dur
}
