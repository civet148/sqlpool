package sqlpool

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"time"
)

type sqlNode struct {
	Priority SqlPriority //优先级
	Event    *sqlEvent   //节点对象
}

type SqlQueue struct {
	Name    string
	Timeout int
	sqlList *sqlList
}

func InstallSqlQueue(config *SqlConfig) (err error) {

	if err = installDatabase(config); err != nil {
		log.Errorf("panic: install database error [%v]", err.Error())
		return
	}
	//创建全局SQL执行通道
	channel = newSqlChannel(config)
	return
}

func newSqlNode(priority SqlPriority, event *sqlEvent) *sqlNode {
	return &sqlNode{
		Priority: priority,
		Event:    event,
	}
}

/*
创建SQL执行队列
strQueueName     队列名称(队列唯一名称)
sqlType  操作数据库类型(SqlType_Any=任意 SqlType_Master=主库 SqlType_Slave=从库)
timeout  超时时间（默认值通过配置指定）
*/
func NewSqlQueue(strQueueName string, timeout int) *SqlQueue {

	if strQueueName == "" {
		log.Errorf("parameter 'strQueueName' is nil")
		return nil
	} else {
		//检查队列名称是否已存在
		if ok := channel.isQueueExist(strQueueName); ok {
			log.Errorf("queue [%v] is exist", strQueueName)
			return nil
		}
	}

	if timeout == 0 {
		timeout = channel.TimeOut()
	}

	q := &SqlQueue{
		Name:    strQueueName,
		Timeout: timeout,
		sqlList: new(sqlList),
	}
	channel.AddQueue(q)
	return q
}

func (q *SqlQueue) GetName() string {
	return q.Name
}

func (q *SqlQueue) String() string {
	return q.GetName()
}

func (q *SqlQueue) GetTimeOut() int {
	return q.Timeout
}

func (q *SqlQueue) Invoke(priority SqlPriority, handler SqlHandler, request Object) Object {

	var invokeTime = time.Now().Format("2006-01-02 15:04:05")
	var timeout = q.GetTimeOut()
	var strQueueName = q.GetName()
	var event = newSqlEvent(handler, request, timeout)

	log.Debugf("sql queue [%v] invoke handler [%+v] request [%+v]", strQueueName, handler, request)

	q.insertEvent(priority, event)

	//等待队列处理结果
	select {
	case r := <-event.receiving: //处理结果通知
		{
			response := r.(*SqlResponse)

			if response.GetError() != nil {
				log.Errorf("sql queue [%+v] method [%+v] event [%+v] invoke response code [%+v] error [%+v]", strQueueName, event.GetMethod(), event, response.GetCode(), response.GetError())
			} else {
				log.Debugf("sql queue [%+v] method [%+v] event [%+v] invoke response code [%+v] sql elapse %+vms invoke elapse %+vms",
					strQueueName, event.GetMethod(), event, response.GetCode(), response.SqlElapse(), response.InvokeElapse())
			}
			return response
		}
	case <-event.expiring: //处理超时通知
		{
			log.Errorf("sql queue [%+v] request [%+v] invoke time [%v], timeout with [%v] seconds", strQueueName, request.String(), invokeTime, timeout)
			return &SqlResponse{code: SqlCode_Timeout, err: fmt.Errorf("sql queue [%+v] request [%+v] invoke timeout", strQueueName, request.String())}
		}
	}
	return nil
}

func (q *SqlQueue) insertEvent(priority SqlPriority, event *sqlEvent) {

	node := newSqlNode(priority, event)
	q.sqlList.insert(node)
	channel.queueNotify(q.GetName())
}
