package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Current Time: ", time.Now().Local())

	//myDelayTask(3*time.Second, handler.DoDelayBiz)
	//mySleepTask(3*time.Second, handler.DoDelayBiz)
	//timerTask(3*time.Second, handler.DoDelayBiz)
	//tickerTask(3*time.Second, handler.DoDelayBiz)
	//timerTaskV2(3*time.Second, handler.DoDelayBiz)
	//tickerTaskV2(3*time.Second, handler.DoDelayBiz)
}

// only second for demo
func myDelayTask(delayTime time.Duration, biz func()) {
	startTime := time.Now().Local().Unix()
	endTime := startTime + int64(delayTime/time.Second)

	for time.Now().Local().Unix() < endTime {
	}
	biz()
}

func mySleepTask(delayTime time.Duration, biz func()) {
	time.Sleep(delayTime)
	biz()
}

// timerTask use timer
func timerTask(delayTime time.Duration, biz func()) {
	select {
	case <-time.After(delayTime):
		biz()
	}
}

// tickerTask use ticker
func tickerTask(delayTime time.Duration, biz func()) {
	select {
	case <-time.Tick(delayTime):
		biz()
	}
}

func timerTaskV2(delayTime time.Duration, biz func()) {
	timer := time.NewTimer(delayTime)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			biz()
		default:
		}
	}

}

func tickerTaskV2(delayTime time.Duration, biz func()) {
	ticker := time.NewTicker(delayTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			biz()
		}
	}

}
