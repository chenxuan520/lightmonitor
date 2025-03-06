package cron

import (
	"github.com/chenxuan520/lightmonitor/internal/notify"
	"log"
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
			log.Printf("notify %s not found", notifyWay)
			continue
		}
		err := way.Send(n.NotifyMsg)
		if err != nil {
			log.Printf("notify %s error: %v", notifyWay, err)
			continue
		}
	}
}
