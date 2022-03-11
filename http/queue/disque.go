package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/zhucheer/orange/cfg"
	"github.com/gomodule/redigo/redis"
	"os"
	"strconv"
	"sync"
	"time"
)

var iDisque Disque

func RegisterDisque(v Disque) {
	if iDisque != nil {
		panic("iDisque is not nil")
	}

	iDisque = v
}

func GetDisque() Disque {
	if iDisque == nil {
		panic("iDisque is nil")
	}

	return iDisque
}

type Disque interface {

	// # 初始化
	Add(int, int) error
	// 获取一个任务
	Pull() (int, error)
	// 完成一个任务
	Success(int) error
	// 重新加入队列
	Rejoin(int) error
	// 程序退出
	Exit()
}

var (
	ErrEOF    = errors.New("服务关闭")
	ErrParams = errors.New("请求参数不合法")
)

type Task struct {
	Id       int
	Class    int
	JoinTime time.Time
	ReadTime time.Time
	Success  bool
}

// lock 锁
var rwMux = &sync.RWMutex{}

type DisqueImpl struct {
	// 存储数据
	tasks map[int]*Task
	// 正在进行的数据
	doing map[int]*Task
	// 队列
	list BinaryHead
	// 服务状态
	status bool
	// 前缀
	prefix string

	// 检查所有服务已经完成启动
	_start sync.WaitGroup

	// 是否开启debug
	debug bool
}

func (d *DisqueImpl) Lock() {
	rwMux.Lock()
}

func (d *DisqueImpl) Unlock() {
	rwMux.Unlock()
}

func (d *DisqueImpl) RLock() {
	rwMux.RLock()
}

func (d *DisqueImpl) RUnlock() {
	rwMux.RUnlock()
}

func (d *DisqueImpl) Exit() {
	// 修改状态
	d.Lock()
	d.status = false
	d.Unlock()

	// 等待一回儿
	log.Info("关闭服务")
	log.Info("先休眠一回儿")
	time.Sleep(time.Second * 10)

	log.Info("等待同步完成")
	GetSync().WaitSync()

	// 检查同步是否完成
	log.Info("等待sync全部完成")
}

// 创建一个队列
func NewQueue() Disque {
	d := new(DisqueImpl)
	d.init()

	d.Disable()

	// 恢复数据到队列
	d.StartRestTask()
	// 超时检查
	d.StartCheckTimeoutCron()

	go func() {
		// 等待其他服务完成启动
		d._start.Wait()
		// 标记服务正常运行
		d.Enable()
	}()

	return d
}

// 初始数据结构
func (d *DisqueImpl) init() {
	d.tasks = make(map[int]*Task)
	d.doing = make(map[int]*Task)
	d.list = NewBinaryHead()

	// 设置前缀
	d.prefix = cfg.Config.GetString("app.prefix")
	if d.prefix == "" {
		d.prefix, _ = os.Hostname()
	}

	d.debug = cfg.Config.GetBool("app.debug")
}

func (d *DisqueImpl) Debug() bool {
	return d.debug
}

func (d *DisqueImpl) StartCheckTimeoutCron() {
	if d.Debug() {
		log.Info("启动检查超时task")
	}
	go func() {
		d.CheckTimeoutCron()
	}()

}
func (d *DisqueImpl) CheckTimeoutCron() {
	tick := time.NewTicker(time.Second * 3)
	for {
		if d.Running() {
			select {
			case <-tick.C:
				if d.Debug() {
					log.Warning("check task timeout cron")
				}
				d.CheckTimeoutTask()
			}
		} else {
			time.Sleep(time.Second / 5)
		}
	}
}

// 超时未完成任务，重新加入队列
func (d *DisqueImpl) CheckTimeoutTask() {
	d.Lock()
	defer d.Unlock()

	t := time.Now()
	for _, task := range d.doing {
		if t.Sub(task.ReadTime).Seconds() > 60 {
			d._Rejoin(task.Id)
		}
	}
}

func (d *DisqueImpl) Add(class, id int) error {

	if id < 1 {
		return ErrParams
	}

	// 检查服务已经关闭
	// 写入tasks
	// 写入list
	// 同步task->hset

	if d.Eof() {
		return ErrEOF
	}

	// lock
	d.Lock()
	defer d.Unlock()

	// 写入存储中
	if _, ok := d.tasks[id]; ok {
		if d.Debug() {
			log.Warning("add a old task, id: %d", id)
		}
		return nil
	}
	task := &Task{
		Id:       id,
		Class:    class,
		JoinTime: time.Now(),
	}
	d.tasks[id] = task

	// 写入队列中
	d.list.Push(class, id)

	// 写入数据库中
	GetSync().AddSync(func() error {
		if d.Debug() {
			log.Warning("add a new task, id: %d,time:%s", id, task.JoinTime.String())
		}
		err := Redis(func(runner redis.Conn) error {
			rwMux.RLock()
			body, _ := json.Marshal(task)
			rwMux.RUnlock()
			_, err := redis.Int(runner.Do("hset", d.redisHkey(), strconv.Itoa(int(id)), body))
			return err
		})
		return err
	})

	return nil
}

func (d *DisqueImpl) redisHkey() string {
	return fmt.Sprintf("%s_data", d.prefix)
}

// 服务已经关闭
func (d *DisqueImpl) Eof() bool {

	return !d.Running()
}

func (d *DisqueImpl) Eject() (int, error) {

	if d.Eof() {
		return 0, ErrEOF
	}

	// list 弹出一个数据
	// 写入doing
	d.Lock()
	defer d.Unlock()

	e, ok := d.list.Pop()
	if !ok {
		if d.Debug() {
			log.Warning("eject a task, list is null")
		}
		return 0, nil
	}

	task := e.Val.(*Task)
	task.ReadTime = time.Now()
	d.doing[task.Id] = task

	return task.Id, nil
}

func (d *DisqueImpl) Pull() (int, error) {

	// 尝试一次
	id, err := d.Eject()
	if err != nil {
		return 0, err
	}
	if id > 0 {
		return id, nil
	}

	ticker := time.NewTimer(time.Second * 60)
	tick := time.NewTicker(time.Second / 10)
	defer ticker.Stop()
	defer tick.Stop()
	for {
		select {
		case <-ticker.C:
			return 0, nil
		case <-tick.C:
			id, err = d.Eject()
			if err != nil {
				return 0, err
			}
			if id > 0 {
				return id, nil
			}
		}
	}
}

// 服务状态，是否运行
func (d *DisqueImpl) Running() bool {
	d.RLock()
	defer d.RUnlock()

	return d.status
}

// 标记服务状态。标记服务正常运行
func (d *DisqueImpl) Enable() {
	log.Info("开启服务")

	d.Lock()
	defer d.Unlock()

	if !d.status {
		d.status = true
	}
}

// 标记服务状态。将服务关闭
func (d *DisqueImpl) Disable() {
	d.Lock()
	defer d.Unlock()

	if d.status {
		d.status = false
	}
}

func (d *DisqueImpl) Success(id int) error {

	// doing 移除
	// 标记成功
	// sync task->hset

	if d.Eof() {
		return ErrEOF
	}

	// 标记任务已经完成
	d.Lock()
	defer d.Unlock()

	// doing 标记移除
	task, ok := d.doing[id]
	if !ok {
		return nil
	}
	delete(d.doing, id)
	// 标记任务完成
	task.Success = true
	// 同步到数据库中
	GetSync().AddSync(func() error {
		if d.Debug() {
			log.Warning("complete a task, id: %d,time:%s", id, task.JoinTime.String())
		}
		err := Redis(func(runner redis.Conn) error {
			rwMux.RLock()
			body, _ := json.Marshal(task)
			rwMux.RUnlock()
			_, err := redis.Int(runner.Do("hset", d.redisHkey(), strconv.Itoa(int(id)), string(body)))
			return err
		})

		return err
	})

	return nil
}

func (d *DisqueImpl) StartRestTask() {

	d._start.Add(1)
	go func() {
		defer d._start.Done()
		// 等待其他服务启动完成
		time.Sleep(time.Second * 5)
		d.ResetTask()
	}()
}

// 启动服务。从hset同步到tasks
func (d *DisqueImpl) ResetTask() {

	// 从hset 读取写入tasks
	// 将未完成的tasks写入list
	log.Info("从hset恢复数据到队列中")
	defer log.Info("从hset恢复数据到队列中完成")
	// lock
	d.Lock()
	defer d.Unlock()

	// 从redis中读取数据。还原到d.tasks中
	Redis(func(runner redis.Conn) error {

		// 从hset读取全部数据
		tasks, err := redis.StringMap(runner.Do("hgetall", d.redisHkey()))
		if err != nil {
			if err == redis.ErrNil {
				return nil
			}
			panic(err)
		}

		// 写入tasks中
		for _, body := range tasks {
			var task = &Task{}
			json.Unmarshal([]byte(body), task)
			d.tasks[task.Id] = task
			if !task.Success {
				log.Info("恢复数据 hset -> task, id: %d", task.Id)
				d.list.Push(task.Class, task)
			}
		}

		return err
	})
}

func (d *DisqueImpl) Rejoin(id int) error {

	// doing 移除
	// 写入list

	if d.Eof() {
		return ErrEOF
	}

	// lock
	d.Lock()
	defer d.Unlock()

	d._Rejoin(id)

	return nil
}
func (d *DisqueImpl) _Rejoin(id int) error {

	if d.Debug() {
		log.Warning("rejoin a task, id: %d", id)
	}
	// doing 移除
	// 写入list

	// doing 移除
	task, ok := d.doing[id]
	if !ok {
		return nil
	}
	delete(d.doing, id)
	// 写入队列中, (class+100)规避垃圾数据。死循环
	d.list.Push(task.Class+100, task)

	return nil
}
