package timewheel

import "time"

//DelayTask 时间轮内部使用
type DelayTask struct {
	Job      func()        // 任务需要执行的Job
	Interval time.Duration // 任务间隔时间

	Circle uint64 // 任务需要在轮盘走多少圈才能执行
	Pos    uint64 // 任务在轮盘的位置

	// todo
	createdTime uint64 // 任务的创建时间 timestamp
	times       int    // 任务需要执行的次数，如果需要一直循环执行，设置成<0的数
	taskName    string // 用来标识task对象，是唯一的
}
