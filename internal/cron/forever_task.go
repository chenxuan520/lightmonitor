package cron

import "math"

type ForeverTask struct {
}

func (f *ForeverTask) Name() string {
	return "ForeverTask"
}

func (f *ForeverTask) NextRunTime() int64 {
	return math.MaxInt64
}

func (f *ForeverTask) Run() {
	return
}
