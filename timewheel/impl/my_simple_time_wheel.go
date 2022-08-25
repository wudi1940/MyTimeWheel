package impl

import (
	"TimeWheel/timewheel"
	"TimeWheel/utils"
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
		Slot: make([][]*timewheel.SimpleDelayTask, size),
	}
}

type MySimpleTimeWheel struct {
	timewheel.BasicTimeWheel
	Slot [][]*timewheel.SimpleDelayTask

	// todo
	taskMap map[string]timewheel.SimpleDelayTask
}

func (tw *MySimpleTimeWheel) Start() {
	ticker := tw.Trigger.(*time.Ticker)
	// init tw
	for i := 0; i < int(tw.Size); i++ {
		tw.Slot[i] = make([]*timewheel.SimpleDelayTask, 0, 10)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				//fmt.Println("Current Time: ", t, " CurPoz: ", tw.CurrentPoz)
				go tw.run()
			}
		}
	}()
}

func (tw *MySimpleTimeWheel) run() {
	oriTaskList := tw.Slot[tw.CurrentPoz]

	if len(oriTaskList) > 0 {
		newTaskList := make([]*timewheel.SimpleDelayTask, 0, len(oriTaskList))
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

func (tw MySimpleTimeWheel) AddTask(task interface{}) {
	iTask := task.(*timewheel.SimpleDelayTask)
	iTask.Circle = uint64(iTask.Interval/tw.Interval) / tw.Size
	iTask.Pos = (tw.CurrentPoz + uint64(iTask.Interval/tw.Interval)) % tw.Size
	tw.Slot[iTask.Pos] = append(tw.Slot[iTask.Pos], iTask)
}

func (tw MySimpleTimeWheel) RemoveTask() {
	//TODO implement me
	panic("implement me")
}
