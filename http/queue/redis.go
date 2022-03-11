package queue

import (
	"gitee.com/zhucheer/orange/database"
	"github.com/gomodule/redigo/redis"
)

type rdFun func(runner redis.Conn) error

func Redis(fn rdFun) error {
	//获取一个 Redis 操作对象，参数是配置中对应的名称
	db, put, err := database.GetRedis("default")
	if err != nil {
		return err
	}

	//连接使用完成后记得将连接放回连接池，否则会造成连接池耗尽大量产生短链接等问题
	defer database.PutConn(put)

	return fn(db)

}
