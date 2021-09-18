package sw

type M map[string]interface{}

type HttpRequestType int

const (
	HttpGet HttpRequestType = iota
	HttpPost
	HttpPut
	HttpPatch
	HttpDelete
	HttpHead
	HttpOption
	HttpNotFound HttpRequestType = -1
)

var httpRequestTypeOption = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTION"}

func (h HttpRequestType) String() (result string) {
	if int(h) == -1 || len(httpRequestTypeOption) <= int(h) {
		result = "NOTFOUND"
	} else {
		result = httpRequestTypeOption[h]
	}
	return
}

var httpRequestMapOption = func() (option map[string]HttpRequestType) {
	option = make(map[string]HttpRequestType)
	for i, curType := range httpRequestTypeOption {
		option[curType] = HttpRequestType(i)
	}
	return
}()

func GetRequestType(requestType string) (result HttpRequestType) {
	if find, ok := httpRequestMapOption[requestType]; ok {
		result = find
	} else {
		result = -1
	}
	return
}

type RouteType int

const (
	//RouteTypeStatic 静态路由
	RouteTypeStatic RouteType = iota
	//RouteTypeActive 动态路由
	RouteTypeActive
)
