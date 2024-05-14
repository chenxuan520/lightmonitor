package service

import (
	"log"
	"sort"
	"time"

	"github.com/chenxuan520/lightmonitor/internal/config"
)

type CronTask struct {
	NextRunTime int64
	Monitor     *Monitor
}

// for sort func
type CronTasks []CronTask

func (c CronTasks) Len() int {
	return len(c)
}
func (c CronTasks) Less(i, j int) bool {
	return c[i].NextRunTime < c[j].NextRunTime
}
func (c CronTasks) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Cron struct {
	NotifyIntervalMinutes int32
	Tasks                 []CronTask
	ConfirmChan           chan struct{}
	SnapshotChan          chan chan struct{ Tasks []CronTask }
}

func NewCron() *Cron {
	c := &Cron{
		NotifyIntervalMinutes: int32(config.GlobalConfig.NotifyWay.NotifyIntervalMinutes),
		Tasks:                 []CronTask{},
		ConfirmChan:           make(chan struct{}),
		SnapshotChan:          make(chan chan struct{ Tasks []CronTask }),
	}
	for _, monitor := range config.GlobalConfig.Monitors {
		m := &Monitor{}
		m.Init(&monitor)
		c.Tasks = append(c.Tasks, CronTask{
			NextRunTime: int64(m.IntervalSeconds) + time.Now().Unix(),
			Monitor:     m,
		})
	}
	return c
}

func (c *Cron) Run() {
	sort.Sort(CronTasks(c.Tasks))
	nextRunTask := &c.Tasks[0]
	nextNotifyTime := time.Now().Unix() + int64(c.NotifyIntervalMinutes*60)
	for {
		select {
		case <-time.After(time.Duration(nextRunTask.NextRunTime-time.Now().Unix()) * time.Second):
			log.Println("INFO: Run task ", nextRunTask.Monitor.URL)
			interval := nextRunTask.Monitor.Check()
			nextRunTask.NextRunTime = interval + time.Now().Unix()
			sort.Sort(CronTasks(c.Tasks))
			nextRunTask = &c.Tasks[0]

		case <-time.After(time.Duration(nextNotifyTime-time.Now().Unix()) * time.Second):
			log.Println("INFO: NotifyIntervalMinutes")
			for i := range c.Tasks {
				if c.Tasks[i].Monitor.Status == OfflineWaitNotify {
					go c.Tasks[i].Monitor.Notify()
				}
			}
			nextNotifyTime = time.Now().Unix() + int64(c.NotifyIntervalMinutes*60)

		case <-c.ConfirmChan:
			log.Println("INFO: Receive confirm signal")
			for i := range c.Tasks {
				if c.Tasks[i].Monitor.Status == OfflineNotified {
					c.Tasks[i].Monitor.Reset()
					c.Tasks[i].NextRunTime = int64(c.Tasks[i].Monitor.IntervalSeconds) + time.Now().Unix()
				}
			}
			sort.Sort(CronTasks(c.Tasks))
			nextRunTask = &c.Tasks[0]

		case snapshotChan := <-c.SnapshotChan:
			log.Println("INFO: Receive snapshot signal")
			tasks := make([]CronTask, len(c.Tasks))
			copy(tasks, c.Tasks)
			snapshotChan <- struct{ Tasks []CronTask }{Tasks: tasks}
		}
	}
}
