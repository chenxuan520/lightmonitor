package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Server struct {
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type Monitor struct {
	URL               string   `json:"url"`
	Method            string   `json:"method"`
	IntervalSeconds   int64    `json:"interval_seconds"`
	MaxOfflineMinutes int64    `json:"max_offline_minutes"`
	Remarks           string   `json:"remarks"`
	Notifications     []string `json:"notifications"`
}

type Email struct {
	Domain    string `json:"domain"`
	Password  string `json:"password"`
	SendEmail string `json:"send_email"`
	RecvEmail string `json:"recv_email"`
}

type Wechat struct {
	SendKey string `json:"send_key"`
}

type Feishu struct {
	WebHook string `json:"web_hook"`
}

type NotifyWay struct {
	NotifyIntervalMinutes int32  `json:"notify_interval_minutes"`
	Email                 Email  `json:"email"`
	Wechat                Wechat `json:"wechat"`
	Feishu                Feishu `json:"feishu"`
}

type Config struct {
	Server    Server    `json:"server"`
	Monitors  []Monitor `json:"monitors"`
	NotifyWay NotifyWay `json:"notify_way"`
}

var (
	GlobalConfig *Config
)

func Init() {
	configFile := "./config/config.json"
	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Println(err)
		configFile = "config.json"
		data, err = ioutil.ReadFile("/config/" + configFile)
		if err != nil {
			log.Println("Read config error!")
			log.Panic(err)
			return
		}
	}

	config := &Config{}

	err = json.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return
	}

	GlobalConfig = config
	log.Println("Config " + configFile + " loaded.")
}

func InitWithPath(configPath string) error {
	data, err := ioutil.ReadFile(configPath)

	if err != nil {
		return err
	}

	config := &Config{}

	err = json.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal config error!")
		return err
	}

	GlobalConfig = config
	log.Println("Config " + configPath + " loaded.")
	return nil
}
