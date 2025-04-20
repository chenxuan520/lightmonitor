package cron

import (
	"log"

	"github.com/chenxuan520/lightmonitor/internal/notify"
)

type NotifyCycleTask struct {
	TaskName   string
	NotifyMsg  notify.NotifyMsg
	RunTime    int64
	CycleTime  int64
	NotifyWays []string
}

func (n *NotifyCycleTask) Name() string {
	return n.TaskName
}

func (n *NotifyCycleTask) NextRunTime() int64 {
	return n.RunTime
}

func (n *NotifyCycleTask) Run() {
	for _, notifyWay := range n.NotifyWays {
		way := notify.GetNotifyByStr(notifyWay)
		if way == nil {
			log.Printf("%s notify %s not found", n.TaskName, notifyWay)
			continue
		}
		err := way.Send(n.NotifyMsg)
		if err != nil {
			log.Printf("%s notify %s error: %v", n.TaskName, notifyWay, err)
			continue
		}
	}
	n.RunTime += n.CycleTime
}

func (n *NotifyCycleTask) IsValid() bool {
	if n.RunTime == 0 {
		return false
	}
	if n.CycleTime == 0 {
		return false
	}
	return true
}
