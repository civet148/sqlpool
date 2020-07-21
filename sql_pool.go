package sqlpool

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/redigogo"
	"github.com/civet148/sqlca"
	"strings"
)

type SqlPool struct {
	db    *sqlca.Engine  //数据库对象
	redis redigogo.Cache //redis缓存对象
}

var pool = &SqlPool{} //全局SQL任务池对象

func installDatabase(pool *SqlPool, config *SqlConfig) (err error) {

	pool.db = sqlca.NewEngine()

	if config.Log.LogPath != "" {
		pool.db.SetLogFile(config.Log.LogPath)
	}

	if strings.ToLower(config.Log.LogLevel) == "debug" {
		pool.db.Debug(true)
	}

	for _, v := range config.Mysql.Masters {
		var options = sqlca.Options{Max: v.Active, Idle: v.Idle}
		pool.db.Open(v.Dsn, &options)
	}
	for _, v := range config.Mysql.Slaves {
		var options = sqlca.Options{Max: v.Active, Idle: v.Idle}
		pool.db.Open(v.Dsn, &options)
	}
	return
}

func installRedis(pool *SqlPool, config *SqlConfig) (err error) {

	c := redigogo.Config{
		Password:       config.Redis.Password,
		Index:          config.Redis.Index,
		MasterHost:     config.Redis.MasterHost,
		ReplicateHosts: config.Redis.ReplicateHosts,
		ConnTimeout:    config.Redis.ConnTimeout,
		ReadTimeout:    config.Redis.ReadTimeout,
		WriteTimeout:   config.Redis.WriteTimeout,
		KeepAlive:      config.Redis.KeepAlive,
		AliveTime:      config.Redis.AliveTime,
	}
	pool.redis = redigogo.NewCache(&c)
	return
}

func InstallSqlPool(config *SqlConfig) (err error) {

	if err = installDatabase(pool, config); err != nil {
		log.Errorf("panic: install database error [%v]", err.Error())
		return
	}
	if err = installRedis(pool, config); err != nil {
		log.Errorf("panic: install redis error [%v]", err.Error())
		return
	}
	//创建全局SQL执行通道
	channel = newSqlChannel(config)
	return
}

func (p *SqlPool) Engine() *sqlca.Engine {
	return p.db
}

func (p *SqlPool) Redis() redigogo.Cache {
	return p.redis
}
