package main

import (
	"TimeWheel/handler"
	"TimeWheel/timewheel"
	"TimeWheel/timewheel/impl"
	"time"
)

func main() {

	wheel := impl.NewMySimpleTimeWheel(1*time.Second, 20)

	task := &timewheel.DelayTask{
		Interval: 3 * time.Second,
		Job: func() {
			handler.DoDelayBiz()
		},
	}

	wheel.Start()

	for {
		time.Sleep(4 * time.Second)
		wheel.AddTask(task)
	}

	select {}

}
