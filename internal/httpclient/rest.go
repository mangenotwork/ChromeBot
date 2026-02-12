package httpclient

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
)

// http 的相关操作

type HttpReq struct {
	Method  string
	Url     string
	Body    []byte
	Header  gt.Header
	Ctype   string
	Cookie  gt.Cookie
	Timeout gt.ReqTimeOut
	Proxy   gt.ProxyUrl
	Stress  int
}

func (req *HttpReq) Do() {
	arg := make([]interface{}, 0)
	if req.Header != nil && len(req.Header) > 0 {
		arg = append(arg, req.Header)
	}
	if req.Cookie != nil && len(req.Cookie) > 0 {
		arg = append(arg, req.Cookie)
	}
	if req.Timeout > 0 {
		arg = append(arg, req.Timeout)
	}
	if len(req.Proxy) > 0 {
		arg = append(arg, req.Proxy)
	}

	var ctx *gt.Context
	var err error

	if req.Stress > 0 {
		// 并发请求
		sum := int64(req.Stress)
		total := sum / 2
		if total < 1 {
			total = 1
		}
		stressObj := gt.NewTestUrl(req.Url, req.Method, sum, int(total))
		stressObj.Run()

	} else {

		switch req.Method {
		case "get":
			ctx, err = gt.Get(req.Url, arg...)

		case "post":
			ctx, err = gt.Post(req.Url, req.Body, req.Ctype, arg...)

		case "put":
			ctx, err = gt.Put(req.Url, req.Body, req.Ctype, arg...)

		case "delete":
			ctx, err = gt.Delete(req.Url, arg...)

		case "options":
			ctx, err = gt.Options(req.Url, arg...)

		case "head":
			// todo 需要gt支持
			ctx, err = gt.Get(req.Url, arg...)

		case "patch":
			// todo 需要gt支持
			ctx, err = gt.Get(req.Url, arg...)

		}
	}

	if err != nil {
		fmt.Println("[ERROR]请求失败err = ", err)
	}
	if ctx != nil {
		fmt.Println("返回code ", ctx.StateCode)
		fmt.Println("返回body ", string(ctx.RespBody))
	}

}
