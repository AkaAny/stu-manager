package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	"stu-manager/config"
	"stu-manager/controller"
	"stu-manager/logger"
	"stu-manager/orm"
)

type Student struct {
	ID   int64
	Name string
}

func main() {
	config.InitConfig("config/config.toml")

	orm.CreateStuGOrm()

	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		var name = context.Param("name")
		var student = Student{ID: 1, Name: name}
		context.JSON(200, student)
	})
	controller.GetRoomController().Init(r)
	controller.GetStudentController().Init(r)

	var listenAddress = getListenAddress()
	logger.Info.Printf("gin listen address:%s", listenAddress)
	err := r.Run(listenAddress)
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
}

func getListenAddress() string {
	var ginConfig = GinConfig{Port: 8080} //默认bind在8080端口
	ginTree := config.GetConfig().GetRootTree().Get("gin").(*toml.Tree)
	err := ginTree.Unmarshal(&ginConfig)
	if err != nil {
		logger.Error.Fatalln(err)
		return ""
	}
	return fmt.Sprintf(":%d", ginConfig.Port)
}
