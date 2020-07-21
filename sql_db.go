package sqlpool

import (
	"github.com/civet148/sqlca"
	"strings"
)

var db *sqlca.Engine //全局数据库对象

func installDatabase(config *SqlConfig) (err error) {

	db = sqlca.NewEngine()

	if config.Log.LogPath != "" {
		db.SetLogFile(config.Log.LogPath)
	}

	if strings.ToLower(config.Log.LogLevel) == "debug" {
		db.Debug(true)
	}

	for _, v := range config.Mysql.Masters {
		var options = sqlca.Options{Max: v.Active, Idle: v.Idle}
		db.Open(v.Dsn, &options)
	}
	for _, v := range config.Mysql.Slaves {
		var options = sqlca.Options{Max: v.Active, Idle: v.Idle}
		db.Open(v.Dsn, &options)
	}
	return
}
