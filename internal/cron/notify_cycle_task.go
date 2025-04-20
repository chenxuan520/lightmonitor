package cron

import (
	"github.com/chenxuan520/lightmonitor/internal/notify"
)

type NotifyCycleTask struct {
	TaskName   string
	NotifyMsg  notify.NotifyMsg
	RunTime    int64
	NotifyWays []string
}
