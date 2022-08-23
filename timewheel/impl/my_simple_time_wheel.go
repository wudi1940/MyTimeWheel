package impl

import (
	"TimeWheel/timewheel"
	"TimeWheel/utils"
	"fmt"
	"time"
)

func NewMySimpleTimeWheel(interval time.Duration, size uint64) *MySimpleTimeWheel {
	return &MySimpleTimeWheel{
		BasicTimeWheel: timewheel.BasicTimeWheel{
			Interval:   interval,
			Size:       size,
			CurrentPoz: utils.TimeToSec(time.Now().Local()) % size,
			Trigger:    time.NewTicker(interval),
		},
		Slot: make([][]*timewheel.DelayTask, size),
	}
}

type MySimpleTimeWheel struct {
	timewheel.BasicTimeWheel
	Slot     [][]*timewheel.DelayTask
	timeUnit time.Duration

	// todo
	taskMap map[string]timewheel.DelayTask
}

func (tw *MySimpleTimeWheel) Start() {
	ticker := tw.Trigger.(*time.Ticker)
	// init tw
	for i := 0; i < int(tw.Size); i++ {
		tw.Slot[i] = make([]*timewheel.DelayTask, 0, 10)
	}

	go func() {
		for {
			select {
			case t := <-ticker.C:

				fmt.Println("Current Time: ", t, " CurPoz: ", tw.CurrentPoz)
				go tw.run()
			}
		}
	}()
}

func (tw *MySimpleTimeWheel) run() {
	oriTaskList := tw.Slot[tw.CurrentPoz]

	if len(oriTaskList) > 0 {
		newTaskList := make([]*timewheel.DelayTask, 0, len(oriTaskList))
		for _, task := range oriTaskList {
			if task.Circle == 0 {
				go task.Job()
				continue
			}

			task.Circle--
			newTaskList = append(newTaskList, task)
		}
		tw.Slot[tw.CurrentPoz] = newTaskList
	}

	tw.CurrentPoz = (tw.CurrentPoz + 1) % tw.Size
}

func (tw MySimpleTimeWheel) AddTask(task *timewheel.DelayTask) {
	task.Circle = uint64(task.Interval/tw.Interval) / tw.Size
	task.Pos = (tw.CurrentPoz + uint64(task.Interval/tw.Interval)) % tw.Size

	fmt.Println("Add New Task !!, Now Time is: ", time.Now().Local(), " CurrentPoz: ", tw.CurrentPoz, " TaskPoz: ", task.Pos)
	tw.Slot[task.Pos] = append(tw.Slot[task.Pos], task)
}

func (tw MySimpleTimeWheel) RemoveTask() {
	//TODO implement me
	panic("implement me")
}
