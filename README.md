# lightmonitor
## Author
-  **chenxuan**
## 功能
- 一个轻量化的监控工具
- 通过配置文件进行设置监控
- 支持通过微信和邮箱两种形式发送告警信息
## 安装
1. `git clone https://github.com/chenxuan520/lightmonitor`
2. `make`(确保安装例如go环境和make工具)
3. 编辑`./config/config.json` 文件,具体规则参考下文
4. `./lightmonitor`(默认配置是 ./config/config.json),此时在`http://127.0.0.1:5200` 默认进入管理界面(端口根据配置)
## 配置文件说明
```json
{
	"server": {
		"port": 5200, // 绑定的端口号
		"password": "example" // 密码,为空就是没有
	},
	"monitors": [
		{
			"url": "https://example.com",// 需要检测的链接
			"method": "GET",// 访问的方法
			"interval_seconds": 10,// 检测的时间间隔
			"max_offline_minutes": 1,// 最长离线时间
			"remarks": "",// 备注
			"notifications": ["wechat","email"]// 通知方式
		}
	],
	"notify_way": {
		"notify_interval_minutes": 60,// 通知检查的间隔
		"email": {
			"domain": "smtp.example.com",// 下面配置参考邮箱提供商
			"password": "example",
			"send_email": "example@example.com",
			"recv_email": "example@example.com"
		},
		"wechat": {// 从 server酱 获取(一个负责接入微信推送的机构,直接搜索即可)
			"send_key": "example"
		}
	}

```
