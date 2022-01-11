package sw

import (
	"fmt"
	"net/http"
	"strings"
)

func NewApp() (route Route) {
	return
}

//parseSyntaxTree 解析路由语法树 性能较好，只会在程序启动时解析一次，后续可以加快匹配速度
func (r Route) parseSyntaxTree() (h HandlerContainer) {
	h = HandlerContainer{
		Handlers: make(map[string]HandlersItem),
		element:  make(map[int]map[HttpRequestType][]ActiveElement),
	}
	for _, group := range r.groups {
		groupLen := len(group.group)
		routeLen := len(group.route)
		if groupLen+routeLen == 0 {
			continue
		}
		r.dealPath(&group.group)
		r.dealPath(&group.route)
		routeParseList, routeParamList, matchList := r.parseRoute(group.group + group.route)
		handlerKey := fmt.Sprintf("%s.%s", "/"+strings.Join(matchList, "/"), group.method.String())
		logInfo(handlerKey)
		h.Handlers[handlerKey] = HandlersItem{
			routeParseList: routeParseList,
			routeParamList: routeParamList,
			matchList:      matchList,
			GroupItem:      group,
		}
		matchListLen := len(matchList)
		if matchListLen == 0 {
			continue
		}
		if _, ok := h.element[matchListLen]; !ok && matchListLen > 0 {
			h.element[matchListLen] = make(map[HttpRequestType][]ActiveElement)
		}
		h.element[matchListLen][group.method] = append(h.element[matchListLen][group.method], ActiveElement{
			FirstElement: matchList[0],
			FullRoute:    handlerKey,
		})
	}
	return
}

func (r Route) parseRoute(key string) (regexpKeys []string, paramKeys []string, matchList []string) {
	matchKeys := strings.Split(key, "/")
	for i, k := range matchKeys {
		if len(k) > 0 {
			tmp := matchKeys[i]
			if k[0:1] == ":" {
				matchKeys[i] = "([a-zA-Z0-9.-]*)"
				paramKeys = append(paramKeys, k[1:])
				tmp = "*"
			}
			regexpKeys = append(regexpKeys, matchKeys[i])
			matchList = append(matchList, tmp)
		}
	}
	return
}

func (r Route) dealPath(path *string) {
	newPath := *path
	pathLen := len(newPath)
	if pathLen > 0 && newPath[:1] != "/" {
		newPath = "/" + newPath
	}
	pathLen = len(newPath)
	if pathLen > 1 && newPath[pathLen-1:pathLen] == "/" {
		newPath = newPath[:pathLen-1]
	}
	path = &newPath
}

func (r Route) Run(addr ...string) (err error) {
	//r.genRouters()
	address := resolveAddress(addr)
	logPrint(address, "SERVER START.")
	handlers := r.parseSyntaxTree()
	err = http.ListenAndServe(address, handlers)
	return
}
