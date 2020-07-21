package sqlpool

import (
	"github.com/civet148/gotools/log"
	"reflect"
)

type TupleContext struct {
	Method string //方法名
}

// request parameter name -> method
var tuples = map[string]TupleContext{}

func getTupleContext(v Object) *TupleContext {

	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	strName := typ.Name()
	if tc, ok := tuples[strName]; ok {
		return &tc
	}
	log.Debugf("request parameter name [%v] tuple context not found", strName)
	return nil
}

func Register(name string, ctx TupleContext) {
	tuples[name] = ctx
}
