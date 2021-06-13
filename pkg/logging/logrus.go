package logging

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jamsa/gin-k8s/pkg/setting"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func Setup(){
	appSetting := setting.AppSetting;
	// TODO: 增加日志级别和输出类型配置信息
	//log.ParseLevel("1");
	//if !strings.Contains(cfg.Output, "stdout") {
	//log.SetOutput(ioutil.Discard)
	//}
	if len(appSetting.LogSavePath) != 0 && len(appSetting.LogSaveName) != 0 {
		writer, _ := rotatelogs.New(
			appSetting.LogSavePath+appSetting.LogSaveName+".%Y%m%d%H%M."+appSetting.LogFileExt,
			rotatelogs.WithLinkName(appSetting.LogSavePath+appSetting.LogSaveName+"."+appSetting.LogFileExt), // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Duration(appSetting.LogMaxAge) * time.Minute),             // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(appSetting.LogRotationTime) * time.Minute), // 日志切割时间间隔
		)
		//pathMap := lfshook.WriterMap{
		//logrus.InfoLevel:  writer,
		//logrus.PanicLevel: writer,
		//}
		/*hook =Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		))*/
		hook := lfshook.NewHook(writer, &log.TextFormatter{})
		//lfshook.NewHook()

		log.AddHook(hook)
	}
}

// Debug output logs at debug level
func Debug(v ...interface{}) {
	log.Debug(v)
}

// Info output logs at info level
func Info(v ...interface{}) {
	log.Info(v)
}

// Warn output logs at warn level
func Warn(v ...interface{}) {
	log.Warn(v)
}

// Error output logs at error level
func Error(v ...interface{}) {
	log.Error(v)
}

// Fatal output logs at fatal level
func Fatal(v ...interface{}) {
	log.Fatal(v)
}


/*
 * Gin中间件函数，记录请求日志
 */
func LogToLogrus() gin.HandlerFunc{
	return func(c *gin.Context){
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := fmt.Sprintf("%6v",endTime.Sub(startTime))

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		log.WithFields(log.Fields{
			"http_status": statusCode,
			"total_time" : latencyTime,
			"ip" : clientIP,
			"method" : reqMethod,
			"uri" : reqUri,
		}).Info("access")
	}
}
