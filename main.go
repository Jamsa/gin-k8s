package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jamsa/gin-k8s/models"
	"github.com/jamsa/gin-k8s/pkg/gredis"
	"github.com/jamsa/gin-k8s/pkg/logging"
	"github.com/jamsa/gin-k8s/pkg/setting"
	"github.com/jamsa/gin-k8s/pkg/util"
	"github.com/jamsa/gin-k8s/routers"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	util.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)
	// 热更新
	// server.ListenAndServe()

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	// endless.DefaultReadTimeOut = readTimeout
	// endless.DefaultWriteTimeOut = writeTimeout
	// endless.DefaultMaxHeaderBytes = maxHeaderBytes
	// server := endless.NewServer(endPoint, routersInit)
	// server.BeforeBegin = func(add string) {
	// 	log.Printf("Actual pid is %d", syscall.Getpid())
	// }

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
