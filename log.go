package snailframe

import  (
	log "github.com/sirupsen/logrus"
)

type logConfig struct {
	Level int64
	Report bool
}

//初始化日志
func initLog(cfg logConfig) {
	log.SetLevel(log.Level(cfg.Level))
	log.SetReportCaller(cfg.Report) //是否开启详细信息

	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true                    // 显示完整时间
	customFormatter.TimestampFormat = "2006-01-02 15:04:05" // 时间格式
	customFormatter.DisableTimestamp = false                // 禁止显示时间
	customFormatter.DisableColors = false                   // 禁止颜色显示

	log.SetFormatter(customFormatter)
	//https://blog.csdn.net/wslyk606/article/details/81670713
}