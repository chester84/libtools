package libtools

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// 通过反射更改结构体的值
func ChangeValueByColName(aExt interface{}, colName string, dstValue interface{}) interface{} {
	v := reflect.ValueOf(aExt)
	v = v.Elem()

	col := v.FieldByName(colName)
	if !col.IsValid() {
		logs.Error("[tools.changeValueByColName] reflect err. aExt: %#v colName: %s", aExt, colName)
		return aExt
	}

	col.Set(reflect.ValueOf(dstValue))
	return aExt
}

func SetByFields(aExt interface{}, colName string, dstValue interface{}) error {
	aa := reflect.ValueOf(aExt).Elem()

	field := FieldByName(aa, colName)

	if field.IsValid() {
		switch field.Type().Kind() {
		case reflect.Slice:
			//logs.Info("dstValue:%v", dstValue)
			v := reflect.Append(field, reflect.ValueOf(dstValue))
			field.Set(v)
		default:
			field.Set(reflect.ValueOf(dstValue))
		}

	} else {
		logs.Error("[SetByFields] field Is not Valid. colName: %v aExt: %v", colName, aExt)
		err := fmt.Errorf("[SetByFields] field Is not Valid. colName: %v aExt: %v", colName, aExt)
		return err
	}
	return nil
}

// 通过反射取得特定字段
func FieldByName(rv reflect.Value, colName string) reflect.Value {
	//logs.Debug("colName:%v ", colName)
	if colName == "" {
		return reflect.Value{}
	}

	index := strings.Index(colName, ".")
	if index != -1 {
		name := colName[0:index]
		//logs.Info("[FieldByName] name:%v", name)

		field := rv.FieldByName(name)

		//logs.Info("[FieldByName] name:%v field:%#v", name, field)
		return FieldByName(field, colName[index+1:])
	}
	return rv.FieldByName(colName)
}
