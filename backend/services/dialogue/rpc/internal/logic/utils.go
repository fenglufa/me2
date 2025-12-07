package logic

import "time"

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
