package service

import (
	"context"
	"log"
	"math"
	"sync"
	"sync/atomic"

	"github.com/chenxuan520/lightmonitor/internal/config"
	"github.com/chenxuan520/lightmonitor/internal/notify"
	utils "github.com/chenxuan520/lightmonitor/internal/utils"
)

type MonitorStatus int8

const (
	Online            MonitorStatus = 0
	OfflineWaitNotify MonitorStatus = 1
	OfflineNotified   MonitorStatus = 2
)

type Monitor struct {
	URL               string   `json:"url"`
	Method            string   `json:"method"`
	IntervalSeconds   int64    `json:"interval_seconds"`
	MaxOfflineMinutes int64    `json:"max_offline_minutes"`
	Remarks           string   `json:"remarks"`
	Notifications     []string `json:"notifications"`

	FailTimes  int64         `json:"fail_times"`
	RunTimes   int64         `json:"run_times"`
	Status     MonitorStatus `json:"status"`
	StatusLock sync.Mutex    `json:"-"`
}

func (m *Monitor) Init(c *config.Monitor) {
	m.URL = c.URL
	m.Method = c.Method
	m.IntervalSeconds = c.IntervalSeconds
	m.MaxOfflineMinutes = c.MaxOfflineMinutes
	m.Remarks = c.Remarks
	m.Notifications = c.Notifications

	m.FailTimes = 0
	m.RunTimes = 0
	m.Status = Online
	m.StatusLock = sync.Mutex{}

	if m.IntervalSeconds == 0 {
		log.Printf("WARNING: Monitor %s interval_seconds is 0, set to 60 seconds\n", m.URL)
		m.IntervalSeconds = 60
	}
	if m.MaxOfflineMinutes == 0 {
		log.Printf("WARNING: Monitor %s max_offline_minutes is 0, set to 1 minutes\n", m.URL)
		m.MaxOfflineMinutes = 1
	}
}

func (m *Monitor) Check() int64 {
	if m.Status == OfflineNotified {
		// 返回会最大的时间,避免频繁检查
		return math.MaxInt32
	}

	go func() {
		_, err := utils.HttpRequest(context.Background(), m.Method, m.URL, nil, nil, nil)

		// 避免长期加锁导致的卡死问题
		m.StatusLock.Lock()
		defer m.StatusLock.Unlock()

		if m.Status == OfflineNotified {
			return
		}

		atomic.AddInt64(&m.RunTimes, 1)
		if err != nil {
			atomic.AddInt64(&m.FailTimes, 1)
			if m.FailTimes*m.IntervalSeconds >= m.MaxOfflineMinutes*60 {
				m.Status = OfflineWaitNotify
			}
			log.Printf("ERROR: Monitor %s check failed: %s,now err time %d\n", m.URL, err, m.FailTimes)
		} else {
			m.Status = Online
			atomic.StoreInt64(&m.FailTimes, 0)
			log.Printf("INFO: Monitor %s check success\n", m.URL)
		}
	}()

	return m.IntervalSeconds
}

func (m *Monitor) Notify() {
	if m.Status != OfflineWaitNotify {
		return
	}

	log.Printf("INFO: Monitor %s notify start\n", m.URL)
	for _, notifyWay := range m.Notifications {
		msg := notify.NotifyMsg{
			Title:   "Monitor Notify",
			Content: m.Remarks + " " + m.URL + " is offline",
		}
		err := notify.SendNotify(notifyWay, msg)
		if err != nil {
			log.Printf("ERROR: Monitor %s notify way %s failed: %s\n", m.URL, notifyWay, err)
		} else {
			log.Printf("INFO: Monitor %s notify way %s success\n", m.URL, notifyWay)
		}
	}

	m.StatusLock.Lock()
	defer m.StatusLock.Unlock()

	m.Status = OfflineNotified
}

func (m *Monitor) Reset() {
	m.StatusLock.Lock()
	defer m.StatusLock.Unlock()

	m.Status = Online
	atomic.StoreInt64(&m.FailTimes, 0)
}
