package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
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

	//new sql queue
	pool := sqlpool.NewSqlQueue("SQL-QUEUE", config.Queue.Timeout)
	if pool == nil {
		panic("new sql queue failed")
	}

	//invoke...
	obj := pool.Invoke(sqlpool.SqlPriority_High, &SqlPoolDAO{}, &SqlRequest{Id: 1})
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

type SqlRequest struct {
	Id int
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

func (dao *SqlPoolDAO) OnSqlProcess(db *sqlca.Engine, request sqlpool.Object) (response sqlpool.Object, err error) {
	log.Infof("SqlPoolDAO -> OnSqlProcess...request [%+v]", request)
	// your database operation code...
	return &SqlResult{Ok: true, LastInsertId: 100}, nil
}
