package sw

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

type GroupItem struct {
	group       string
	route       string
	method      HttpRequestType
	middlewares []HandlerFunc
	handler     HandlerFunc
}

type Route struct {
	groupMiddlewares []HandlerFunc
	middlewares      []HandlerFunc
	groups           []GroupItem
	group            string
}

func (r *Route) GET(route string, handlerFunc HandlerFunc) {
	r.appendResource(route, handlerFunc, HttpGet)
}

func (r *Route) POST(route string, handlerFunc HandlerFunc) {
	r.appendResource(route, handlerFunc, HttpPost)
}

func (r *Route) PUT(route string, handlerFunc HandlerFunc) {
	r.appendResource(route, handlerFunc, HttpPut)
}

func (r *Route) PATCH(route string, handlerFunc HandlerFunc) {
	r.appendResource(route, handlerFunc, HttpPatch)
}

func (r *Route) DELETE(route string, handlerFunc HandlerFunc) {
	r.appendResource(route, handlerFunc, HttpDelete)
}

func (r *Route) HEAD(route string, handlerFunc HandlerFunc) {
	r.appendResource(route, handlerFunc, HttpHead)
}

func (r *Route) OPTION(route string, handlerFunc HandlerFunc) {
	r.appendResource(route, handlerFunc, HttpOption)
}

//VUE 支持VUE
func (r *Route) VUE(staticPath string, fs *embed.FS, temp ...interface{}) {
	route := "/"
	handlerFunc := HandlerFunc(func(this *This) {
		var staticFile string
		var byteData []byte
		if path := this.GetRequest().URL.Path; path == "/" {
			staticFile = fmt.Sprintf("%s/index.html", staticPath)
		} else {
			staticFile = fmt.Sprintf("%s%s", staticPath, path)
		}
		for _, s := range []string{"./", "../"} {
			staticFile = strings.ReplaceAll(staticFile, s, "")
		}
		byteData, _ = fs.ReadFile(staticFile)
		// 根据文件扩展名设置响应头Content-Type信息
		var contentType string
		switch split := strings.Split(staticFile, "."); split[len(split)-1] {
		case "html", "htm", "xhtml":
			contentType = "text/html"
			staticCache := this.Get("html_"+staticFile, []byte{})
			staticCacheByte := staticCache.([]byte)
			if tempLen := len(temp); tempLen > 0 && len(staticCacheByte) == 0 {
				c, _ := template.New("member").Parse(string(byteData))
				var buf bytes.Buffer
				c.Execute(&buf, temp[0])
				byteData = buf.Bytes()
				this.Set("html_"+staticFile, byteData)
			} else if tempLen > 0 {
				byteData = staticCacheByte
			}
		case "css":
			contentType = "text/css"
		case "js":
			contentType = "text/javascript"
			if gzData, err := fs.ReadFile(staticFile + ".gz"); err == nil {
				byteData = gzData
				this.GetResponse().Header().Set("Content-Encoding", "gzip")
				this.GetResponse().Header().Set("Vary", "Accept-Encoding")
				this.GetResponse().Header().Set("Content-Length", fmt.Sprintf("%v", len(gzData)))
			}
		case "gif", "png", "jpg", "jpeg", "bmp", "ico":
			contentType = "image/*"
		default:
			contentType = "text/plain"
		}
		this.GetResponse().Header().Set("content-type", contentType)
		this.GetResponse().Write(byteData)
	})
	r.appendResource("/*", handlerFunc, HttpGet)
	r.appendResource(route+"assets/:*", handlerFunc, HttpGet)
}

func (r *Route) appendResource(route string, handlerFunc HandlerFunc, method HttpRequestType) {
	if r.groupMiddlewares != nil {
		r.middlewares = append(r.groupMiddlewares, r.middlewares...)
	}
	r.groups = append(r.groups, GroupItem{
		group:       r.group,
		route:       route,
		method:      method,
		middlewares: r.middlewares,
		handler:     handlerFunc,
	})
	r.clear()
}

func (r *Route) clear() {
	r.middlewares = nil
}

func (r *Route) Use(midFunc HandlerFunc) *Route {
	r.middlewares = append(r.middlewares, midFunc)
	return r
}

func (r *Route) UseGroupMid(midFunc HandlerFunc) *Route {
	r.groupMiddlewares = append(r.groupMiddlewares, midFunc)
	return r
}

func (r *Route) Group(name string) *Route {
	r.group = name
	r.groupMiddlewares = nil
	return r
}
