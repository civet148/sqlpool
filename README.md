# sqlpool

## 调用方式

* 初始化SQL任务池
```go
if err := sqlpool.InstallSqlPool(config); err != nil {
    panic("install sql pool failed")
}
```

* 新建队列（队列名唯一）
```go
pool := sqlpool.NewSqlQueue("SQL-QUEUE-20200721190532", config.Queue.Timeout)
if pool == nil {
    panic("new sql queue failed")
}
pool.Invoke(...) //调用Invoke方法等待处理结果
```

* 完整测试案例

```go
package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlpool"
)

const (
	CONFIG_FILE_PATH = "test.toml"
)

func main() {

	var err error
	var config = &sqlpool.SqlConfig{}
	if _, err = toml.DecodeFile(CONFIG_FILE_PATH, config); err != nil {

		log.Errorf("decode toml file [%v] error [%v]", CONFIG_FILE_PATH, err.Error())
		return
	}

	//install sql pool ...
	if err = sqlpool.InstallSqlPool(config); err != nil {
		panic("install sql pool failed")
	}
	log.Debugf("config [%+v]", config)

	//new sql queue
	pool := sqlpool.NewSqlQueue("SQL-QUEUE", config.Queue.Timeout)
	if pool == nil {
		panic("new sql queue failed")
	}

	//invoke...
	obj := pool.Invoke(sqlpool.SqlPriority_High, &SqlPoolDAO{}, &SqlRequest{User: &User{Name: "john", Phone: "8613022223333"}})
	response := obj.(*sqlpool.SqlResponse)
	if response.OK() {
		log.Infof("response ok, result [%+v]", response.Object())
	} else {
		log.Infof("response code [%v] error [%+v]", response.GetCode(), response.GetError())
	}
	select {}
}

type SqlPoolDAO struct {
}

type SqlResult struct {
	Ok           bool
	LastInsertId int64
}

type User struct {
	Id    int32  `db:"id"`
	Name  string `db:"name"`
	Phone string `db:"phone"`
	Sex   int    `db:"sex"`
}

type SqlRequest struct {
	User *User
}

func (r *SqlResult) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *SqlRequest) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (dao *SqlPoolDAO) String() string {
	data, _ := json.Marshal(dao)
	return string(data)
}

func (dao *SqlPoolDAO) OnSqlProcess(pool *sqlpool.SqlPool, request sqlpool.Object) (response sqlpool.Object, err error) {
	log.Infof("SqlPoolDAO -> OnSqlProcess...request [%+v]", request)
	// your database operation code...

	req := request.(*SqlRequest)

	var lastInsertId int64
	if lastInsertId, err = pool.Engine().Model(req.User).Table("users").Insert(); err != nil {
		log.Errorf("SQL insert error [%+v]", err.Error())
		return &SqlResult{Ok: false, LastInsertId: 0}, nil
	}
	return &SqlResult{Ok: true, LastInsertId: lastInsertId}, nil
}


```

## 配置文件

```toml
# 日志配置
[log]
logPath = "test.log"
logLevel = "debug"

[[mysql.masters]]
# 主数据库(读写)
name = "master01"
dsn = "mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4"
active = 100
idle = 2

[[mysql.slaves]]
# 备份数据库(只读)
name = "slave01"
dsn = "mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4&slave=true"
active = 100
idle = 2

[queue]
# 默认协程数(最好小于等于mysql最大连接数)
routines = 100
# 默认超时（秒）
timeout = 5
# 默认容量
cap=50
# 输出调试队列数据
debug = true

[redis]
# redis auth password
password = ""
# redis db index (0~15)
index = 0
# redis master host
masterHost = "127.0.0.1"
# redis replicate hosts
replicateHosts = []
# default connect timeout (ms)
connTimeout = 1000
# default connect timeout (ms)
readTimeout = 1000
# default connect timeout (ms)
writeTimeout = 1000
# default keep alive 30s
keepAlive = 30
# default alive time 60s
aliveTime = 60
```

	