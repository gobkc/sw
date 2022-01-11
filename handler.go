package sw

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type HandlerFunc func(*This)

type HandlerContainer struct {
	Handlers map[string]HandlersItem
	element  map[int]map[HttpRequestType][]ActiveElement //int保存元素个数
}

type ActiveElement struct {
	FirstElement string
	FullRoute    string
}

type HandlersItem struct {
	routeParseList []string
	routeParamList []string
	matchList      []string
	GroupItem
}

func (h HandlerContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	uri := r.URL.Path + "." + r.Method
	var curGroupItem GroupItem
	var routeParse []string
	var routeType RouteType
	//1.先查找静态路由 效率最高，建议不要用动态路由。仅仅是方便而已。算法再好的动态路由都有性能损耗
	if staticRoute, ok := h.Handlers[uri]; ok {
		curGroupItem = staticRoute.GroupItem
		routeType = RouteTypeStatic
	} else {
		routeType = RouteTypeActive
		curGroupItem, routeParse = h.findActiveRoute(uri, r.Method)
	}
	var this = &This{
		w:          w,
		r:          r,
		route:      curGroupItem.route,
		method:     curGroupItem.method.String(),
		routeType:  routeType,
		routeParse: routeParse,
	}
	if curGroupItem.middlewares != nil {
		for _, middleware := range curGroupItem.middlewares {
			if this.abort {
				return
			}
			middleware(this)
		}
	}
	if curGroupItem.route != "" {
		if this.abort {
			return
		}
		if curGroupItem.handler != nil {
			curGroupItem.handler(this)
		}
	}
	cost := time.Since(start)
	logDefault(curGroupItem.method.String(), curGroupItem.route, "->", r.RequestURI, " COST TIME:", fmt.Sprintf("%.5f", cost.Seconds()), "S")
}

//findActiveRoute 查找动态路由
func (h HandlerContainer) findActiveRoute(uri string, requestType string) (curGroupItem GroupItem, routeParse []string) {
	//1.查找动态路由
	matchUris := strings.Split(uri, "/")
	//2. 匹配条件1 查找是否长度符合的route
	matchIndex := len(matchUris)
	if matchIndex > 1 {
		matchIndex--
	}
	findRequestType := GetRequestType(requestType)
	var findFlag = false
	if findElement, ok := h.element[matchIndex]; ok && findRequestType != -1 {
		//3. 匹配条件2 查找请求方法是否一致
		if matchList, matchOK := findElement[findRequestType]; matchOK {
			findProcess := h.FindFirstEle(matchUris[1], matchList)
			findProcess(func(msg ActiveElement) {
				if elementObj, findObjOK := h.Handlers[msg.FullRoute]; findObjOK {
					//5.开始正则匹配
					re := "^/" + strings.Join(elementObj.routeParseList, "/") + "." + requestType + "$"
					parse := regexp.MustCompile(re)
					findArr := parse.FindAllStringSubmatch(uri, -1)
					if len(findArr) > 0 {
						curGroupItem = elementObj.GroupItem
						routeParse = findArr[0]
						findFlag = true
					}
				}
			}, func() {
				findFlag = false
			})
		}
	}
	//如果路由找不到，开始找全局路由和404路由
	if findFlag == false {
		if findPublic, ok := h.Handlers["/*."+requestType]; ok {
			curGroupItem = findPublic.GroupItem
		} else {
			if find404, ok404 := h.Handlers["/404."+requestType]; ok404 {
				curGroupItem = find404.GroupItem
			}
		}
	}
	return
}

func (h HandlerContainer) FindFirstEle(findEle string, elementList []ActiveElement) (result func(success func(msg ActiveElement), fail func())) {
	result = func(success func(msg ActiveElement), fail func()) {
		fail()
	}
	for i, element := range elementList {
		if findEle == element.FirstElement || element.FirstElement == "*" {
			result = func(success func(msg ActiveElement), fail func()) {
				success(elementList[i])
			}
			break
		}
	}
	return
}
