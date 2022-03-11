package main

import (
	"gitee.com/zhucheer/orange/app"
	"queue/route"
)

func main() {

	router := &route.Route{}

	app.AppStart(router)
}
