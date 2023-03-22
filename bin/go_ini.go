package bin

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

var structName string

func loadInfo(fileName string, data interface{}) (fileRead []byte, sTypeof reflect.Type, sValueof reflect.Value, err error) {

	//使用data参数反射获取 类型
	sTypeof = reflect.TypeOf(data)
	sValueof = reflect.ValueOf(data)
	fmt.Println(sTypeof, sTypeof.Kind())
	//检查是否使用了指针类型
	if sTypeof.Kind() != reflect.Ptr {
		err = fmt.Errorf("%v not ptr", err)
		return
	}
	//检查是否属于指针类型的结构体
	if sTypeof.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("%v not ptr  type Struct ", err)
		return
	}

	//
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	//
	return file, sTypeof, sValueof, err
}

//read  conf.ini  file
func FileRead(fileRead []byte, sTypeof reflect.Type, sValueof reflect.Value) (err error) {

	//把ini 文件存在数组中
	lines := strings.Split(string(fileRead), "\r\n")

	//把ini 文件一行一行读 行号：idx(o行开始) ，内容：line， 遍历
	for idx, line := range lines {

		//首尾空格去除
		line = strings.TrimSpace(line)

		//中间的行如果没有内容，跳过
		if len(line) == 0 {
			continue
		}

		//每行以 ; # 符号开头的，跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		//如果是[开头的话
		if strings.HasPrefix(line, "[") {
			//检查开头和结尾是否存在 [ ]
			if line[0] != '[' || line[len(line)-1] != ']' {
				err = fmt.Errorf("The %d line has a syntax error", idx+1)
				return
			}

			//满足[]开头与结尾，就把节点的名字提取
			sectionname := strings.TrimSpace(line[1 : len(line)-1])

			//当[ ]中间内容为空是，提示语法错误
			if len(sectionname) == 0 {
				err = fmt.Errorf("line %d", idx+1)
				return
			}

			//结构体反射,把结构体字段遍历
			for i := 0; i < sTypeof.Elem().NumField(); i++ {
				field := sTypeof.Elem().Field(i)

				//如果节点的名字等于结构体字段的tag标签
				if sectionname == field.Tag.Get("ini") {
					structName = field.Name
					fmt.Printf("find %s The corresponding nested struct  %s\n", sectionname, structName)
				}
			}
		} else {
			//检查 键值对
			if strings.Index(line, "=") == -1 || strings.HasPrefix(line, "=") {
				err = fmt.Errorf("The %d line has a syntax error", idx+1)
				return
			}
			index := strings.Index(line, "=")
			key := strings.TrimSpace(line[:index])
			value := strings.TrimSpace(line[index+1:])
			v := sValueof
			sValue := v.Elem().FieldByName(structName)
			sType := sValue.Type()

			//
			if sType.Kind() != reflect.Struct {
				err = fmt.Errorf("not struct %s", structName)
				return
			}

			var fieldName string
			var fileType reflect.StructField
			for i := 0; i < sValue.NumField(); i++ {
				field := sType.Field(i)
				fileType = field
				if field.Tag.Get("ini") == key {
					fieldName = field.Name
					break
				}
			}
			if len(fieldName) == 0 {
				continue
			}
			//
			fileObj := sValue.FieldByName(fieldName)
			fmt.Println(fieldName, fileType.Type.Kind())

			//ini文件中对值的判断
			switch fileType.Type.Kind() {

			//字符串
			case reflect.String:
				fileObj.SetString(value)

				//整数
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var valueint int64
				valueint, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					err = fmt.Errorf("The %d line has a syntax error", idx+1)
					return
				}
				fileObj.SetInt(valueint)

				//布尔
			case reflect.Bool:
				var valueBool bool
				valueBool, err = strconv.ParseBool(value)
				if err != nil {
					err = fmt.Errorf("The %d line has a syntax error", idx+1)
					return
				}
				fileObj.SetBool(valueBool)

				//浮点型
			case reflect.Float32, reflect.Float64:
				var valueFloat float64
				valueFloat, err = strconv.ParseFloat(value, 64)
				if err != nil {
					err = fmt.Errorf("The %d line has a syntax error", idx+1)
					return
				}
				fileObj.SetFloat(valueFloat)
			}
		}
	}
	return
}
