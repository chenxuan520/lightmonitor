package cron

import (
	"log"

	"github.com/chenxuan520/lightmonitor/internal/notify"
)

type NotifyOnceTask struct {
	TaskName   string
	NotifyMsg  notify.NotifyMsg
	RunTime    int64
	NotifyWays []string
}

func (n *NotifyOnceTask) Name() string {
	return n.TaskName
}

func (n *NotifyOnceTask) NextRunTime() int64 {
	return n.RunTime
}

func (n *NotifyOnceTask) Run() {
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
}

func (n *NotifyOnceTask) IsValid() bool {
	if n.RunTime == 0 {
		return false
	}
	return true
}
