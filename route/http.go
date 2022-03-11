package route

import (
	"gitee.com/zhucheer/orange/app"
	"gitee.com/zhucheer/orange/logger"
	"queue/http/controller"
	"queue/http/middleware"
	"queue/http/queue"
)

type Route struct {
}

const Version = "v1.1.1"

func (s *Route) ServeMux() {

	queue := app.NewRouter("/queue", middleware.NewRecover())
	queue.GET("/add", controller.QueueAdd)
	queue.GET("/pull", controller.QueuePull)
	queue.GET("/success", controller.QueueSuccess)
	queue.GET("/rejoin", controller.QueueRejoin)
	queue.GET("/version", func(ctx *app.Context) error {
		return ctx.ToString("version:" + Version)
	})

}

// Register 服务注册器, 可以对业务服务进行初始化调用
func (s *Route) Register() {
	queue.RegisterDisque(queue.NewQueue())
	queue.ResisterSync(queue.NewSync())

	app.AppDefer(func() {
		logger.Warning("程序进入退出")
		queue.GetDisque().Exit()
		logger.Warning("程序进入退出完成")
	})
}
