package sqlpool

import (
	"encoding/json"
	"github.com/civet148/sqlca"
	"math/rand"
	"reflect"
	"time"
)

type Object interface {
	String() string
}

type SqlHandler interface {
	OnSqlProcess(db *sqlca.Engine, request Object) (response Object, err error)
}

type sqlEvent struct {
	RandomId    int64       //随机数
	Timeout     int         //超时时间(秒)
	CreatedAt   int64       //创建时间戳
	CreatedNano int64       //请求时间(纳秒)
	Request     Object      //请求数据
	Method      string      //请求名称
	receiving   chan Object //处理结果接收通道
	expiring    chan bool   //超时通知
	handler     SqlHandler  //处理对象
}

func newSqlEvent(handler SqlHandler, request Object, timeout int) *sqlEvent {

	var method string
	var now = time.Now()
	var createdAt = now.Unix()
	var createdNano = now.UnixNano()
	ctx := getTupleContext(request)
	if ctx != nil {
		method = ctx.Method
	} else {
		typ := reflect.TypeOf(request)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		method = typ.Name()
	}
	return &sqlEvent{
		Timeout:     timeout,
		CreatedAt:   createdAt,
		CreatedNano: createdNano,
		RandomId:    rand.Int63(),
		receiving:   make(chan Object, 1),
		expiring:    make(chan bool, 1),
		Method:      method,
		Request:     request,
		handler:     handler,
	}
}

func (e *sqlEvent) GetMethod() string {
	return e.Method
}

func (e *sqlEvent) String() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *sqlEvent) GetRequest() Object {
	return e.Request
}

func (e *sqlEvent) GetCreatedNano() int64 {
	return e.CreatedNano
}

func (e *sqlEvent) GetCreatedAt() int64 {
	return e.CreatedAt
}

func (e *sqlEvent) GetTimeOut() int {
	return e.Timeout
}

func (e *sqlEvent) GetRandomId() int64 {
	return e.RandomId
}

func (e *sqlEvent) IsTimeOut() bool {
	now64 := time.Now().Unix()
	if now64-(e.CreatedAt+int64(e.Timeout)) > 0 {
		return true
	}
	return false
}

func (e *sqlEvent) NotifyExpiring() {
	e.expiring <- true
}

func (e *sqlEvent) NotifyResponse(response *SqlResponse) {
	e.receiving <- response
}
