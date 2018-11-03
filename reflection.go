package hyutil

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//ObjFill はmap[string]stringを変換してmodelに展開します。
func ObjFill(model interface{}, row map[string]string) {

	tp := reflect.TypeOf(model)
	val := reflect.ValueOf(model)

	if tp.Kind() == reflect.Ptr {
		val = val.Elem()
		tp = tp.Elem()
	}

	if tp.Kind() != reflect.Struct {
		log.Println("ObjFill:タイプ取得不正")
		return
	}

	for i := 0; i < tp.NumField(); i++ {

		field := tp.Field(i)
		colname := field.Tag.Get("json")

		if colname == "" {
			colname = strings.ToLower(CamelToSnake(field.Name))
		}

		valStr, ok := row[colname]

		if ok {
			dest := val.FieldByName(field.Name)
			convData(&dest, valStr)
		}

	}

}

func convData(dest *reflect.Value, valStr string) {

	if !dest.CanSet() {
		return
	}

	//フィールドのTypeによって文字列から変換
	switch dest.Kind() {
	case reflect.String:
		dest.SetString(valStr)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.Atoi(valStr)
		dest.SetInt(int64(i))
	case reflect.Float32:
		f, _ := strconv.ParseFloat(valStr, 32)
		dest.SetFloat(float64(f))
	case reflect.Float64:
		f, _ := strconv.ParseFloat(valStr, 64)
		dest.SetFloat(float64(f))
	case reflect.Bool:
		b := false
		if valStr == "1" {
			b = true
		}
		dest.SetBool(b)
	case reflect.Struct:
		switch dest.Interface().(type) {
		case DateTime:

			if valStr != "" {
				t, _ := time.Parse(dbTimeFormat, valStr)

				set := reflect.ValueOf(DateTime{
					Time: &t,
				})

				dest.Set(set)
			}

		}

	default:
		log.Println("no case " + dest.Kind().String())
	}

}
