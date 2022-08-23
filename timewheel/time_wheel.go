package timewheel

import "time"

type TimeWheel interface {
	Start()

	AddTask(task *DelayTask)
	RemoveTask()
}

type BasicTimeWheel struct {
	Interval   time.Duration // 轮盘精度
	Size       uint64        // 轮盘大小
	CurrentPoz uint64        // 当前位置

	Trigger interface{}     // 触发时间轮推进
	Slot    [][]interface{} // 轮盘实体，每个element存储Task列表
}
