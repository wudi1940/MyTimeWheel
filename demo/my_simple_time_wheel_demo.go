package main

import (
	"TimeWheel/handler"
	"TimeWheel/timewheel"
	"TimeWheel/timewheel/impl"
	"fmt"
	"time"
)

func main() {

	wheel := impl.NewMySimpleTimeWheel(1*time.Millisecond, 20)
	wheel.Start()

	demonstrateDelayTaskS(wheel)

	select {}

}

func demonstrateDelayTaskS(tw *impl.MySimpleTimeWheel) {
	task := &timewheel.SimpleDelayTask{
		BaseDelayTask: timewheel.BaseDelayTask{
			Interval: 3 * time.Second,
			Job: func() {
				handler.DoDelayBiz()
			},
		},
	}

	for {
		tw.AddTask(task)
		fmt.Println("Add New Task !!, Now Time is: ", time.Now().Local())
		time.Sleep(4 * time.Second)
	}
}
