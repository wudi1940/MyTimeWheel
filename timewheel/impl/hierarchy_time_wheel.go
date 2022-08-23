package impl

import (
	"TimeWheel/timewheel"
	"TimeWheel/utils"
	"sync/atomic"
	"time"
	"unsafe"
)

/*
HierarchyTimeWheel 层级时间轮
 与普通时间轮不同的地方在于，对于超出时间轮范围(Size * Interval)的任务，不通过circle记录任务的位置
而是将该任务升级存储到上一级的时间轮中，上一级的时间轮Interval = 次级时间轮整个范围
*/
type HierarchyTimeWheel struct {
	timewheel.BasicTimeWheel

	Slot               [][]*timewheel.SimpleDelayTask
	parentWheelPointer unsafe.Pointer // *HierarchyTimeWheel 指向上级时间轮的指针
}

func NewHierarchyTimeWheel(interval time.Duration, size uint64, trigger interface{}) *HierarchyTimeWheel {
	return &HierarchyTimeWheel{
		BasicTimeWheel: timewheel.BasicTimeWheel{
			Interval:   interval,
			Size:       size,
			CurrentPoz: utils.TimeToSec(time.Now().Local()) % size,
			Trigger:    trigger,
		},
		Slot: make([][]*timewheel.SimpleDelayTask, size),
	}
}

func (tw *HierarchyTimeWheel) Start() {
	// init tw
	for i := 0; i < int(tw.Size); i++ {
		tw.Slot[i] = make([]*timewheel.SimpleDelayTask, 0, 10)
	}

	if tw.Trigger != nil {
		ticker := tw.Trigger.(*time.Ticker)
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
}

func (tw *HierarchyTimeWheel) run() {
	// 判断是否是最底层tw，只有最底层tw存在trigger
	if tw.Trigger != nil {
		oriTaskList := tw.Slot[tw.CurrentPoz]

		// HierarchyTimeWheel no need to check circle, out bound task will be stored in parentTW
		if len(oriTaskList) > 0 {
			for _, task := range oriTaskList {
				go task.Job()

			}
			tw.Slot[tw.CurrentPoz] = make([]*timewheel.SimpleDelayTask, 0, 10)
		}
	}

	// when current tw run a round, trigger parentTW move
	if tw.CurrentPoz+1 == tw.Size && tw.parentWheelPointer != nil {
		tw.parentTWFillInTW()
	}
	tw.CurrentPoz = (tw.CurrentPoz + 1) % tw.Size
}

func (tw *HierarchyTimeWheel) AddTask(task interface{}) {
	iTask := task.(*timewheel.SimpleDelayTask)

	// 如果任务间隔时间超过当前tw一轮时长，需要升级到parentTW中存储
	if uint64(iTask.Interval) >= uint64(tw.Interval)*tw.Size {
		var parentTW *HierarchyTimeWheel
		// 判断是否已存在parentTW，否则创建
		if tw.getParentTW() != nil {
			parentTW = tw.getParentTW()
		} else {
			parentTW = NewHierarchyTimeWheel(time.Duration(uint64(tw.Interval)*tw.Size), tw.Size, nil)
			parentTW.Start()
			atomic.CompareAndSwapPointer(&tw.parentWheelPointer, nil, unsafe.Pointer(parentTW))
		}
		// 调用parentTW的AddTask
		parentTW.AddTask(task)
		return
	}

	iTask.Pos = (tw.CurrentPoz + uint64(iTask.Interval/tw.Interval)) % tw.Size
	tw.Slot[iTask.Pos] = append(tw.Slot[iTask.Pos], iTask)
}

func (tw *HierarchyTimeWheel) RemoveTask() {
	//TODO implement me
	panic("implement me")
}

func (tw *HierarchyTimeWheel) parentTWFillInTW() {
	parentTW := tw.getParentTW()
	// 如果上级TW存在
	if parentTW != nil {
		// 当前位置存在任务，整个list对应sonTW，再分配到sonTW上
		if len(parentTW.Slot[parentTW.CurrentPoz]) != 0 {
			for _, task := range parentTW.Slot[parentTW.CurrentPoz] {
				tmpTask := &timewheel.SimpleDelayTask{
					BaseDelayTask: timewheel.BaseDelayTask{
						Interval: task.Interval % parentTW.Interval,
						Job:      task.Job,
					},
				}
				tw.AddTask(tmpTask)
			}
		}
		// 清空parentTW当前位置的任务，以及全部降级存储到sonTW上
		parentTW.Slot[parentTW.CurrentPoz] = make([]*timewheel.SimpleDelayTask, 0, 10)
	}

	// 降级后触发parentTW前进一个interval
	go parentTW.run()
}

func (tw *HierarchyTimeWheel) getParentTW() *HierarchyTimeWheel {
	pointer := atomic.LoadPointer(&tw.parentWheelPointer)
	// 如果上级TW存在
	if pointer != nil {
		return (*HierarchyTimeWheel)(pointer)
	}
	return nil
}
