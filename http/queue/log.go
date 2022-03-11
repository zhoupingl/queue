package queue

import (
	"gitee.com/zhucheer/orange/logger"
)

var log *logger.Logger

func init() {
	log = logger.New(logger.INFO, "text", "", 500)
	log.Format.TimeFormat = "2006-01-02 15:04:05"
	log.Format.ContentFormat = "{{.Color}}{{.LevelString}} [{{.Time}}]  {{.Message}}{{.KvJson}} {{.ColorClear}}\n"
}
