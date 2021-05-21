package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"itwork-bbs365/controllers"
	"itwork-bbs365/model"
	"itwork-bbs365/pkg/common"
	"itwork-bbs365/pkg/config"
	"itwork-bbs365/scheduler"
)

var configFile = flag.String("config", "./itwork-bbs365.yaml", "配置文件路径")

func init() {
	flag.Parse()

	// 初始化配置
	conf := config.Init(*configFile)

	// gorm配置
	gormConf := &gorm.Config{}

	// 初始化日志
	if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
		if conf.ShowSql {
			gormConf.Logger = logger.New(log.New(file, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel:      logger.Info,
			})
		}
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Error(err)
	}

	// 连接数据库
	if err := simple.OpenDB(conf.MySqlUrl, gormConf, 10, 20, model.Models...); err != nil {
		logrus.Error(err)
	}
}

func main() {
	if common.IsProd() {
		// 开启定时任务
		scheduler.Start()
	}
	controllers.Router()
}
