package sw

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

//This callback封装
type This struct {
	w          http.ResponseWriter
	r          *http.Request
	route      string
	method     string
	routeType  RouteType
	routeParse []string
	abort      bool
	cache      sync.Map
}

//Json 输出JSON结果
func (t *This) Json(httpStatus int, data interface{}) {
	t.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	t.w.WriteHeader(httpStatus)
	byteData, _ := json.Marshal(data)
	_, err := t.w.Write(byteData)
	if err != nil {
		logPrint(err.Error())
	}
}

//Html 输出HTML结果
func (t *This) Html(httpStatus int, data interface{}) {
	t.w.Header().Set("Content-Type", "text/html;charset=utf-8\n")
	t.w.WriteHeader(httpStatus)
	byteData := []byte(fmt.Sprintf("%v",data))
	_, err := t.w.Write(byteData)
	if err != nil {
		logPrint(err.Error())
	}
}

//WriteHeaderStatus 写响应头状态码
func (t *This) WriteHeaderStatus(httpStatus int) {
	t.w.WriteHeader(httpStatus)
}

//GetResponse 获取响应对象
func (t *This) GetResponse() (w http.ResponseWriter) {
	w = t.w
	return
}

//GetRequest 获取请求对象
func (t *This) GetRequest() (r *http.Request) {
	r = t.r
	return
}

//GetRequestHeaders 获取请求头列表
func (t *This) GetRequestHeaders() (headerList M) {
	headerList = make(M, len(t.r.Header))
	for hKey, hValue := range t.r.Header {
		headerList[hKey] = strings.Join(hValue, ";")
	}
	return
}

//GetRequestHeader 获取指定请求头
func (t *This) GetRequestHeader(headerKey, defaultValue string) (value string) {
	if value = t.r.Header.Get(headerKey); value == "" {
		value = defaultValue
	}
	return
}

//SetRequestHeader 设置请求头
func (t *This) SetRequestHeader(headerKey, headerValue string) {
	t.r.Header.Set(headerKey, headerValue)
	return
}

//SetResponseHeader 设置响应头
func (t *This) SetResponseHeader(headerKey, headerValue string) {
	t.w.Header().Set(headerKey, headerValue)
}

//GetResponseHeader 获取指定相应头
func (t *This) GetResponseHeader(headerKey, defaultValue string) (value string) {
	if value = t.w.Header().Get(headerKey); value == "" {
		value = defaultValue
	}
	return
}

//GetResponseHeaders 获取响应头列表
func (t *This) GetResponseHeaders() (headerList M) {
	headers := t.w.Header()
	headerList = make(M, len(headers))
	for hKey, hValue := range headers {
		headerList[hKey] = strings.Join(hValue, ";")
	}
	return
}

//Params 获取路由参数列表
func (t *This) Params() (m M) {
	m = make(M)
	routeParams := strings.Split(t.route, "/")
	routeParams = func(arr []string) (newArr []string) {
		for _, s := range arr {
			if len(s) >= 2 {
				if s[:1] == ":" {
					newArr = append(newArr, s[1:])
				}
			}
		}
		return
	}(routeParams)
	routeParamLen := len(routeParams)
	for i, s := range t.routeParse {
		if i > 0 && i <= routeParamLen {
			m[routeParams[i-1]] = s
		}
	}
	return
}

//Param 获取路由参数
func (t *This) Param(key, def string) (find string) {
	routeParams := strings.Split(t.route, "/")
	routeParams = func(arr []string) (newArr []string) {
		for _, s := range arr {
			if len(s) >= 2 {
				if s[:1] == ":" {
					newArr = append(newArr, s)
				}
			}
		}
		return
	}(routeParams)
	for i, param := range routeParams {
		if len(param) < 2 {
			continue
		}
		if findRoute := param[1:]; param[:1] == ":" && findRoute == key {
			if len(t.routeParse) >= i {
				find = t.routeParse[i+1]
			}
		}
	}
	if find == "" {
		find = def
	}
	return
}

//Path 获取GET参数
func (t *This) Path(key, def string) (find string) {
	find = t.r.URL.Query().Get(key)
	if find == "" {
		find = def
	}
	return
}

//PathList 获取GET参数列表
func (t *This) PathList() (m M) {
	list := t.r.URL.Query()
	m = make(M, len(list))
	for i, row := range list {
		if rowLen := len(row); rowLen == 1 {
			m[i] = row[0]
		} else if rowLen > 1 {
			m[i] = row
		}
	}
	return
}

//Abort 中止
func (t *This) Abort() {
	t.abort = true
}

//Set 设置
func (t *This) Set(name string, value interface{}) {
	t.cache.Store(name, value)
}

//Get 获取值
func (t *This) Get(name string, def interface{}) (v interface{}) {
	if readValue, ok := t.cache.Load(name); ok {
		v = readValue
		return
	}
	v = def
	return
}

//GetInt 获取int值
func (t *This) GetInt(name string, def int) (v int) {
	if readValue, ok := t.cache.Load(name); ok {
		if v, ok = readValue.(int); ok {
			return
		}
	}
	v = def
	return
}

//GetFloat 获取float值
func (t *This) GetFloat(name string, def float64) (v float64) {
	if readValue, ok := t.cache.Load(name); ok {
		if v, ok = readValue.(float64); ok {
			return
		}
	}
	v = def
	return
}

//GetString 获取string值
func (t *This) GetString(name string, def string) (v string) {
	if readValue, ok := t.cache.Load(name); ok {
		if v, ok = readValue.(string); ok {
			return
		}
	}
	v = def
	return
}

//GetBody 获取消息体
func (t *This) GetBody() (v []byte) {
	v, _ = ioutil.ReadAll(t.r.Body)
	return
}

//ParseJson 从消息体中解析json
func (t *This) ParseJson(parseTo interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadAll(t.r.Body); err != nil {
		return
	}
	if err = json.Unmarshal(data, parseTo); err != nil {
		return
	}
	return
}

//ParseForm 从消息体中解析form参数
func (t *This) ParseForm(parseTo interface{}) (err error) {
	switch formType := t.GetFormType(); formType {
	case "multipart/form-data":
		if err = t.r.ParseMultipartForm(1048576); err != nil {
			return
		}
	case "form-urlencoded":
		if err = t.r.ParseForm(); err != nil {
			return
		}
	default:
		err = errors.New("error form type")
		return
	}
	formParams := t.r.Form
	var m = make(map[string]string, len(formParams))
	for curKey, curRow := range formParams {
		switch curRowLen := len(curRow); curRowLen {
		case 0:
			m[curKey] = ""
		default:
			m[curKey] = curRow[0]
		}
	}
	switch pType := reflect.TypeOf(parseTo).Kind(); pType {
	case reflect.Ptr:
		dataHandler := reflect.ValueOf(parseTo).Elem()
		typeHandler := reflect.TypeOf(parseTo).Elem()
		dataNumField := dataHandler.NumField()
		for i := 0; i < dataNumField; i++ {
			tag, ok := typeHandler.Field(i).Tag.Lookup("json")
			if !ok {
				tag = dataHandler.Type().Field(i).Name
			}
			obj := dataHandler.Field(i)
			t.setFormValue(&obj, m[tag])
		}
	}
	return
}

func (t *This) setFormValue(obj *reflect.Value, newV string) {
	oldValue := obj.Interface()
	switch pType := reflect.TypeOf(oldValue).Kind(); pType {
	case reflect.Uint8:
	case reflect.Int8:
	case reflect.Uint:
	case reflect.Int64:
	case reflect.Uint64:
	case reflect.Int32:
	case reflect.Uint32:
	case reflect.Int:
		tmp, _ := strconv.ParseInt(newV, 10, 64)
		obj.SetInt(tmp)
	case reflect.String:
		obj.SetString(newV)
	case reflect.Float32:
	case reflect.Float64:
		tmp, _ := strconv.ParseFloat(newV, 64)
		obj.SetFloat(tmp)
	case reflect.Bool:
		if newV == "true" {
			obj.SetBool(true)
		} else {
			obj.SetBool(false)
		}
	default:
		return
	}
}

func (t *This) GetFormType() (formType string) {
	formType = "multipart/form-data"
	if findType, ok := t.r.Header["Content-Type"]; ok {
		findTypeReal := ""
		if len(findType) > 0 {
			findTypeReal = findType[0]
		}
		if inPos := strings.Index(findTypeReal, "multipart"); inPos == -1 {
			formType = "form-urlencoded"
		}
	}
	return
}
