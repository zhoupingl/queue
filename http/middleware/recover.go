package middleware

import (
	"fmt"
	"gitee.com/zhucheer/orange/app"
	"gitee.com/zhucheer/orange/logger"
)

type HashTable map[string]interface{}

type Recover struct {
}

func NewRecover() *Recover {
	return &Recover{}
}

func (w *Recover) Func() app.MiddlewareFunc {
	return func(next app.HandlerFunc) app.HandlerFunc {

		return func(c *app.Context) (err error) {
			defer func() {
				if e := recover(); e != nil {
					logger.Error("recover panic: %+v", e)
					c.SetResponseStatus(500)
					err = c.ToJson(HashTable{
						"errno":  400,
						"errmsg": fmt.Sprintf("%v", e),
					})
				}
			}()

			return next(c)
		}

	}
}
