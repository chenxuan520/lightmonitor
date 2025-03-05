package cron

import (
	"errors"
	"log"
	"sort"
	"time"
)

type CronTask interface {
	// 任务名称唯一标识符,相同的任务名称会被认为是同一个任务
	Name() string
	// GetNextRunTime 获取下次运行时间, 如果是周期任务, 则返回下次运行时间, 如果是一次性任务, 则返回-1
	NextRunTime() int64
	// Run 执行任务
	Run()
}

// for sort func
type CronTasks []CronTask

func (c CronTasks) Len() int {
	return len(c)
}

func (c CronTasks) Less(i, j int) bool {
	return c[i].NextRunTime() < c[j].NextRunTime()
}

func (c CronTasks) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CronStatus int8

const (
	Inited CronStatus = 0
	Runing CronStatus = 1
)

type Cron struct {
	Status         CronStatus
	Tasks          []CronTask
	logger         *log.Logger
	AddTaskChan    chan CronTask
	DeleteTaskChan chan string
	SnapshotChan   chan chan struct{ Tasks []CronTask }
}

func NewCron(logger *log.Logger) *Cron {
	if logger == nil {
		logger = log.Default()
	}
	return &Cron{
		Status:         Inited,
		Tasks:          []CronTask{},
		logger:         logger,
		AddTaskChan:    make(chan CronTask),
		DeleteTaskChan: make(chan string),
		SnapshotChan:   make(chan chan struct{ Tasks []CronTask }),
	}
}

func (c *Cron) AddTask(task CronTask) error {
	if task.NextRunTime() == -1 || task.NextRunTime() < time.Now().Unix() {
		return errors.New("task next run time is invalid")
	}
	c.logger.Printf("INFO: Add task %s", task.Name())
	if c.Status == Runing {
		c.AddTaskChan <- task
	} else {
		c.Tasks = append(c.Tasks, task)
	}
	return nil
}

func (c *Cron) Run() {
	c.Status = Runing

	// add forever task
	c.AddTask(&ForeverTask{})

	sort.Sort(CronTasks(c.Tasks))
	nextRunTask := c.Tasks[0]
	for {
		select {
		case <-time.After(time.Duration(nextRunTask.NextRunTime() - time.Now().Unix())):

			c.logger.Printf("INFO: Run task %s", nextRunTask.Name())
			nextRunTask.Run()
			if nextRunTask.NextRunTime() == -1 || nextRunTask.NextRunTime() < time.Now().Unix() {
				c.Tasks = c.Tasks[1:]
			}

			sort.Sort(CronTasks(c.Tasks))
			if len(c.Tasks) == 0 {
				// 理论上不会执行到这里
				return
			}
			nextRunTask = c.Tasks[0]

		case newTask := <-c.AddTaskChan:
			if newTask.NextRunTime() == -1 || newTask.NextRunTime() < time.Now().Unix() {
				continue
			}

			// 去重,replace
			for _, task := range c.Tasks {
				if task.Name() == newTask.Name() {
					task = newTask
					break
				}
			}

			c.Tasks = append(c.Tasks, newTask)
			sort.Sort(CronTasks(c.Tasks))
			nextRunTask = c.Tasks[0]

		case taskName := <-c.DeleteTaskChan:
			for i, task := range c.Tasks {
				if task.Name() == taskName {
					c.Tasks = append(c.Tasks[:i], c.Tasks[i+1:]...)
					break
				}
			}
			sort.Sort(CronTasks(c.Tasks))
			nextRunTask = c.Tasks[0]

		case snapshotChan := <-c.SnapshotChan:
			tasks := make([]CronTask, len(c.Tasks))
			copy(tasks, c.Tasks)
			snapshotChan <- struct{ Tasks []CronTask }{Tasks: tasks}
		}
	}
}
