package job

import (
	"time"
)

func TimeStamp() int64 {
	return time.Now().Unix()
}

func Order() int64 {
	return time.NowNano().Unix()
}
