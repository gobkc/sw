package sw

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
)

//BasicAuth basic auth 中间件
func BasicAuth(account, password string) HandlerFunc {
	return func(this *This) {
		rHeader := this.GetRequestHeader("Authorization", "")
		var eType = fmt.Sprintf("%s:%s", account, password)
		var esEncode = base64.StdEncoding.EncodeToString([]byte(eType))
		baseAuthString := fmt.Sprintf("Basic %s", esEncode)
		if rHeader == "" || rHeader != baseAuthString {
			this.SetResponseHeader("WWW-Authenticate", "Basic realm="+strconv.Quote("Authorization Required"))
			this.WriteHeaderStatus(http.StatusUnauthorized)
			this.Abort()
			return
		}
	}
}

//Cors 跨域 中间件 这里只是通用的方式。建议自己实现
func Cors() HandlerFunc {
	return func(this *This) {
		method := this.GetRequest().Method
		this.SetResponseHeader("Access-Control-Allow-Origin", "*")
		this.SetResponseHeader("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		this.SetResponseHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT,PATCH,HEAD")
		this.SetResponseHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization, Token")
		this.SetResponseHeader("Access-Control-Allow-Credentials", "true")
		/*放行所有OPTIONS方法*/
		if method == "OPTIONS" {
			this.WriteHeaderStatus(http.StatusNoContent)
			this.Abort()
			return
		}
	}
}
