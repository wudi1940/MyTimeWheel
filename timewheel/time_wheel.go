package timewheel

import "time"

type TimeWheel interface {
	// Start 初始化tw，并运行
	Start()

	// AddTask 添加延迟任务到tw
	AddTask(task interface{})
	// RemoveTask 删除tw中的延迟任务
	RemoveTask()
}

type BasicTimeWheel struct {
	Interval   time.Duration // 轮盘精度
	Size       uint64        // 轮盘大小
	CurrentPoz uint64        // 当前位置

	Trigger interface{}     // 触发时间轮推进
	Slot    [][]interface{} // 轮盘实体，每个element存储Task列表

	// todo 待加入功能，支持延迟任务管理
	taskMap map[string]BaseDelayTask
}
