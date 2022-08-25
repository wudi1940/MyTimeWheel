package main

import (
	"TimeWheel/handler"
	"TimeWheel/timewheel"
	"TimeWheel/timewheel/impl"
	"fmt"
	"time"
)

const timeUnit = time.Millisecond

func main() {

	ticker := time.NewTicker(timeUnit)
	wheel := impl.NewHierarchyTimeWheel(1*timeUnit, 20, ticker)
	wheel.Start()

	demonstrateDelayTaskH(wheel)

	select {}

}

func demonstrateDelayTaskH(tw *impl.HierarchyTimeWheel) {
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