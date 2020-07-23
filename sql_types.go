package sqlpool

type SqlType int

const (
	SqlType_Any    SqlType = 0 //操作主库或从库
	SqlType_Master SqlType = 1 //操作主库
	SqlType_Slave  SqlType = 2 //操作从库
)

func (t SqlType) GoString() string {
	return t.String()
}

func (t SqlType) String() string {
	switch t {
	case SqlType_Any:
		return "SqlType_Any"
	case SqlType_Master:
		return "SqlType_Master"
	case SqlType_Slave:
		return "SqlType_Slave"
	}
	return "SqlType_Unknown"
}

type SqlCode int

const (
	SqlCode_OK             SqlCode = 0 //成功
	SqlCode_Timeout        SqlCode = 1 //超时
	SqlCode_Error_Database SqlCode = 2 //数据库错误
	SqlCode_Error_Handler  SqlCode = 3 //程序执行错误
	SqlCode_Error_Closed   SqlCode = 4 //队列已关闭
)

func (c SqlCode) GoString() string {
	return c.String()
}

func (c SqlCode) String() string {
	switch c {
	case SqlCode_OK:
		return "SqlCode_OK"
	case SqlCode_Timeout:
		return "SqlCode_Timeout"
	case SqlCode_Error_Database:
		return "SqlCode_Error_Database"
	case SqlCode_Error_Handler:
		return "SqlCode_Error_Handler"
	case SqlCode_Error_Closed:
		return "SqlCode_Error_Closed"
	}
	return "SqlCode_Unknown"
}

type SqlPriority int

const (
	SqlPriority_Low  SqlPriority = 0 //优先级：低
	SqlPriority_Mid  SqlPriority = 1 //优先级：中
	SqlPriority_High SqlPriority = 2 //优先级：高
)

func (c SqlPriority) GoString() string {
	return c.String()
}

func (c SqlPriority) String() string {
	switch c {
	case SqlPriority_Low:
		return "SqlPriority_Low"
	case SqlPriority_Mid:
		return "SqlPriority_Mid"
	case SqlPriority_High:
		return "SqlPriority_High"
	}
	return "SqlPriority_Unknown"
}
