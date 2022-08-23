package handler

import (
	"fmt"
	"time"
)

func DoDelayBiz() {
	fmt.Println("Doing Delay Biz..., Now Time is: ", time.Now().Local(), "\n")
}
