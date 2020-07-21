package sqlpool

import "github.com/civet148/gotools/log"

func InstallSqlPool(config *SqlConfig) (err error) {

	if err = installDatabase(config); err != nil {
		log.Errorf("panic: install database error [%v]", err.Error())
		return
	}
	//创建全局SQL执行通道
	channel = newSqlChannel(config)
	return
}
