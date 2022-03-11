package controller

import (
	"gitee.com/zhucheer/orange/app"
	"queue/http/queue"
	"strconv"
)

type HashTable map[string]interface{}

// 添加
func QueueAdd(c *app.Context) error {

	_id := c.Request().FormValue("id")
	_class := c.Request().FormValue("class")
	id, _ := strconv.Atoi(_id)
	class, _ := strconv.Atoi(_class)
	err := queue.GetDisque().Add(class, id)
	if err != nil {
		c.SetResponseStatus(400)
		return c.ToJson(HashTable{
			"errno":  400,
			"errmsg": err.Error(),
		})
	}

	return c.ToString("ok")
}

// 获取一个
func QueuePull(c *app.Context) error {

	id, err := queue.GetDisque().Pull()
	if err != nil {
		c.SetResponseStatus(400)
		return c.ToJson(HashTable{
			"errno":  400,
			"errmsg": err.Error(),
		})
	}

	return c.ToString(strconv.Itoa(int(id)))
}

// 标记成功
func QueueSuccess(c *app.Context) error {

	_id := c.Request().FormValue("id")
	id, _ := strconv.Atoi(_id)
	err := queue.GetDisque().Success(id)
	if err != nil {
		c.SetResponseStatus(400)
		return c.ToJson(HashTable{
			"errno":  400,
			"errmsg": err.Error(),
		})
	}

	return c.ToString("ok")
}

// 标记失败
func QueueRejoin(c *app.Context) error {

	_id := c.Request().FormValue("id")
	id, _ := strconv.Atoi(_id)
	err := queue.GetDisque().Rejoin(id)
	if err != nil {
		c.SetResponseStatus(400)
		return c.ToJson(HashTable{
			"errno":  400,
			"errmsg": err.Error(),
		})
	}

	return c.ToString("ok")
}
