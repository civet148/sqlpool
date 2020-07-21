package sqlpool

import "fmt"

type SqlResponse struct {
	code         SqlCode //错误码: 0=成功 其他表示失败
	err          error   //错误描述
	object       Object  //返回的数据
	caller       string  //请求名称
	sqlElapse    int64   //SQL handler处理方法执行耗时(毫秒)
	invokeElapse int64   //从调用Invoke到得到返回数据执行耗时(毫秒)
}

//可选参数:
//	args[0] sql execute elapse (int64) SQL handler处理方法执行耗时(毫秒)
//	args[1] invoke execute elapse (int64) 从调用Invoke到得到返回数据执行耗时(毫秒)
func newSqlResponseOk(object Object, args ...interface{}) *SqlResponse {
	var sqlElapse, invokeElapse int64

	if len(args) == 2 {
		if v, ok := args[0].(int64); ok {
			sqlElapse = v
		}
		if v, ok := args[1].(int64); ok {
			invokeElapse = v
		}
	}

	return &SqlResponse{
		code:         SqlCode_OK,
		err:          nil,
		object:       object,
		caller:       object.String(),
		sqlElapse:    sqlElapse,
		invokeElapse: invokeElapse,
	}
}

func newSqlResponseError(code SqlCode, err error) *SqlResponse {
	return &SqlResponse{
		code:   code,
		err:    err,
		object: nil,
	}
}

func (s *SqlResponse) OK() bool {
	return s.code == SqlCode_OK
}

func (s *SqlResponse) GetCode() SqlCode {
	return s.code
}

func (s *SqlResponse) GetError() error {
	return s.err
}

func (s *SqlResponse) Object() Object {
	return s.object
}

func (s *SqlResponse) String() (text string) {

	text += fmt.Sprintf("{")
	text += fmt.Sprintf("code:  %+v ", s.code)
	text += fmt.Sprintf("err:  %+v ", s.err)
	text += fmt.Sprintf("object:  %+v ", s.object)
	text += fmt.Sprintf("caller:  %+v ", s.caller)
	text += fmt.Sprintf("sqlElapse:  %+v ", s.sqlElapse)
	text += fmt.Sprintf("invokeElapse:  %+v ", s.invokeElapse)
	text += fmt.Sprintf("}")
	return
}

func (s *SqlResponse) SqlElapse() int64 {
	return s.sqlElapse
}

func (s *SqlResponse) InvokeElapse() int64 {
	return s.invokeElapse
}
