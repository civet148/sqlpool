package sqlpool

import (
	"github.com/civet148/gotools/log"
	"math/rand"
	"sync"
	"time"
)

const (
	DEFAULT_SQL_QUEUE_TIMEOUT      = 3     //超时时间(单位：秒)
	DEFAULT_SQL_QUEUE_CAP          = 100   //容量
	DEFAULT_SQL_QUEUE_ROUTINES     = 2     //默认一个队列协程数
	DEFAULT_SQL_CHANNEL_NOTIFY_CAP = 10000 //通知通道容量
)

type sqlChannel struct {
	config    *SqlConfig           //配置信息
	name      string               //队列名称
	timeout   int                  //超时时间（秒）
	cap       int                  //队列容量
	routines  int                  //协程数量
	receiving chan *sqlEvent       //事件接收通道
	closing   chan bool            //关闭消息通道
	notifying chan string          //消息通知通道
	isClosed  bool                 //是否已关闭
	locker    sync.RWMutex         //互斥锁
	queues    map[string]*SqlQueue //事件队列
}

var channel *sqlChannel //全局数据库处理通道

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
创建SQL执行队列
timeout  超时时间（默认3秒）
cap      容量(默认100）
routines 处理协程数(默认2)
*/
func newSqlChannel(config *SqlConfig) *sqlChannel {

	var timeout, capabilities, routines int

	routines = config.Queue.Routines
	timeout = config.Queue.Timeout
	capabilities = config.Queue.Cap

	if routines == 0 {
		routines = DEFAULT_SQL_QUEUE_ROUTINES
	}
	if timeout == 0 {
		timeout = DEFAULT_SQL_QUEUE_TIMEOUT
	}
	if capabilities == 0 {
		capabilities = DEFAULT_SQL_QUEUE_CAP
	}

	c := &sqlChannel{
		timeout:   timeout,
		cap:       capabilities,
		routines:  routines,
		config:    config,
		queues:    make(map[string]*SqlQueue, capabilities),
		receiving: make(chan *sqlEvent, capabilities),
		closing:   make(chan bool, 1),
		notifying: make(chan string, DEFAULT_SQL_CHANNEL_NOTIFY_CAP),
	}
	c.start()
	return c
}

func (c *sqlChannel) GetName() string {
	return c.name
}

func (c *sqlChannel) String() string {
	return c.GetName()
}

func (c *sqlChannel) GetCap() int {
	return c.cap
}

func (c *sqlChannel) Close() {
	c.lock()
	defer c.unlock()
	c.isClosed = true
	c.closing <- c.isClosed
	close(c.receiving)
	close(c.closing)
}

func (c *sqlChannel) IsClosed() bool {
	return c.isClosed
}

func (c *sqlChannel) Routines() int {
	return c.routines
}

func (c *sqlChannel) TimeOut() int {
	return c.timeout
}

func (c *sqlChannel) lock() {
	c.locker.Lock()
}

func (c *sqlChannel) unlock() {
	c.locker.Unlock()
}

func (c *sqlChannel) rlock() {
	c.locker.RLock()
}

func (c *sqlChannel) runlock() {
	c.locker.RUnlock()
}

func (c *sqlChannel) addQueue(q *SqlQueue) (ok bool) {
	c.lock()
	defer c.unlock()
	c.queues[q.GetName()] = q
	return
}

func (c *sqlChannel) isQueueExist(strQueueName string) (ok bool) {
	c.rlock()
	defer c.runlock()
	if _, ok := c.queues[strQueueName]; ok {
		return true
	}
	return
}

func (c *sqlChannel) start() {
	//启动协程接收SQL队列消息
	for i := 0; i < c.routines; i++ {
		go c.runLoop()
	}
}

func (c *sqlChannel) runLoop() {

	for {
		select {
		case strQueueName := <-c.notifying: //队列插入成功通知(开始从指定队列取数据进行处理)
			{
				if c.IsClosed() {
					log.Warnf("sql channel [%v] is closed", c.GetName())
					return
				}
				c.onSqlNotify(strQueueName)
			}
		case event := <-c.receiving: //SQL事件处理
			{
				c.onSqlEvent(event)
			}
		case <-c.closing:
			{
				log.Infof("sql channel [%v] closing...", c.GetName())
				break
			}
		}
	}
}

func (c *sqlChannel) AddQueue(q *SqlQueue) {
	c.addQueue(q)
	log.Debugf("queue [%v] add to channel map ok", q.GetName())
}

func (c *sqlChannel) onSqlNotify(strQueueName string) {
	c.rlock()
	defer c.runlock()

	if v, ok := c.queues[strQueueName]; ok {
		if c.isQueueDebug() {
			//调试模式，输出队列信息(通过配置项开启/关闭)
			v.sqlList.print(strQueueName)
		}
		node := v.sqlList.front()
		event := node.Event
		if event.IsTimeOut() {
			event.NotifyExpiring() //SQL事件已超时，通知Invoke方法
		} else {
			c.queueReceiving(event) //处理SQL事件
		}
	} else {
		log.Errorf("queue name [%v] not found", strQueueName)
	}
}

func (c *sqlChannel) onSqlEvent(event *sqlEvent) {
	log.Debugf("OnSqlEvent request [%+v]", event.Request)

	var response *SqlResponse

	//处理结果发送完毕关闭通道
	defer close(event.receiving)
	//计时
	var nano1 = time.Now().UnixNano()

	//调用处理方法得到结果
	object, err := event.handler.OnSqlProcess(db, event.Request)

	var nano2 = time.Now().UnixNano()

	sqlElapseMS := (nano2 - nano1) / 1000000                     //SQL处理方法执行耗时（毫秒）
	invokeElapseMS := (nano2 - event.GetCreatedNano()) / 1000000 //invoke调用整体耗时（毫秒）
	//将结果返回给调用者
	if err != nil {
		response = newSqlResponseError(SqlCode_Error_Handler, err) //方法执行失败
	} else {
		response = newSqlResponseOk(object, sqlElapseMS, invokeElapseMS) //执行成功
	}
	event.NotifyResponse(response)
}

func (c *sqlChannel) isQueueDebug() (ok bool) {
	if c.config != nil && c.config.Queue.Debug {
		return true
	}
	return
}

func (c *sqlChannel) queueNotify(strQueueName string) {
	c.notifying <- strQueueName
}

func (c *sqlChannel) queueReceiving(event *sqlEvent) {
	c.receiving <- event
}
