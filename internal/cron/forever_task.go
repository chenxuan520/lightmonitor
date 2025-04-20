package cron

import (
	"time"
)

type ForeverTask struct {
}

func (f *ForeverTask) Name() string {
	return "ForeverTask"
}

func (f *ForeverTask) NextRunTime() int64 {
	// 100 年后
	return time.Now().AddDate(100, 0, 0).Unix()
}

func (f *ForeverTask) Run() {
	return
}

func (f *ForeverTask) IsValid() bool {
	return true
}
