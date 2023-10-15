package MYHTTP

import (
	"fmt"
	"net"
	"strings"
)

type Connect struct {
	PORT string //规定端口
}

// 定义context函数类型
type Handler func(context *Context)

// 定义创建json函数里的 map类型
type H map[string]any

var Handlers = make(map[string]Handler)

type HTTPStruct struct {
	requestMethod       string            //请求方法
	requestHost         string            //请求域名
	requestUrlPath      string            //请求url路径
	getQuery            map[string]string //get请求参数
	postQuery           map[string]string
	responseCode        string //响应状态码
	responseContentType string //响应内容类型
	responseDate        string //响应时间
	//responseContentLength string   			//响应内容长度
	responseBody string   //响应内容
	Engine       *Context //规定行为的引擎 Engine
}

func New() HTTPStruct {
	return HTTPStruct{
		requestMethod:       "",
		requestHost:         "",
		requestUrlPath:      "",
		getQuery:            nil,
		postQuery:           nil,
		responseCode:        "",
		responseContentType: "",
		responseDate:        "",
		responseBody:        "",
		Engine:              nil,
	}
}

func (h *HTTPStruct) MainConnetc(port string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		fmt.Println("启动失败，端口或被占用")
		return
	}
	fmt.Printf("HTTP服务端口启动在: %s\n", port)
	// 监听成功后，等待客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}
		// 启动一个goroutine处理客户端连接
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		message := string(buf)
		fmt.Println(message)
		h := Parser(message, conn)

		//fmt.Println(h.requestUrlPath)

		//urlpath := getrequestcompleteurlpath(h.responseBody)

		//
		//if h.requestMethod == "GET"{
		//
		//
		//}

		switch h.requestMethod {
		case "GET":
			//num := strings.Index(urlpath, "GET")
			//handler, ok := Handlers[getrequesturlpath(h.responseBody)]
			//if ok || num != -1 {
			//	handler(h.Engine)
			//	h.Engine.Conn.Close()
			//}

			handler, ok := Handlers[getrequesturlpath(getrequestcompleteurlpath(h.responseBody))]
			if ok {
				handler(h.Engine)
				h.Engine.Conn.Close()
			}
			return
		case "POST":

			handler, ok := Handlers[getrequesturlpath(getrequestcompleteurlpath(h.responseBody))]
			if ok {
				handler(h.Engine)
				h.Engine.Conn.Close()
			}
			return
		default:
			fmt.Println("不正确的tcp链接")
			return
		}

	}

}

func Parser(requestbody string, conn net.Conn) *HTTPStruct {

	return &HTTPStruct{
		requestHost:    getrequesthost(requestbody),
		requestUrlPath: getrequesturlpath(getrequestcompleteurlpath(requestbody)),
		requestMethod:  getrequestmethod(requestbody),
		getQuery:       getquery(getrequestcompleteurlpath(requestbody)),
		//postQuery:      postquery(requestbody),
		responseBody: requestbody,

		Engine: CreateEngine(conn,
			getrequesthost(requestbody),
			getrequesturlpath(getrequestcompleteurlpath(requestbody)),
			getquery(getrequestcompleteurlpath(requestbody)),
			postquery(requestbody),
		),
	}

}

func CreateEngine(conn net.Conn, requesthost string, requesturlpath string, query map[string]string, postquery map[string]string) *Context {
	return &Context{
		Conn:           conn,
		RequestHost:    requesthost,
		RequestUrlPath: requesturlpath,
		getQuery:       query,
		PostQuery:      postquery,
		WriteBody:      "HTTP/1.1 200\nContent-Type: text/plain; charset=utf-8",
	}
}

func getrequesthost(requestbody string) string {

	index01 := strings.Index(requestbody, "\nHost: ")
	index02 := strings.Index(requestbody, "\nConnection")
	return requestbody[index01+7 : index02]
}

func getrequestmethod(requestbody string) string {
	index01 := strings.Index(requestbody, " ")
	return requestbody[0:index01]
}

func getrequesturlpath(completeurlpath string) string {
	//取出来url参数符号 ?

	index01 := strings.Index(completeurlpath, "/")
	index02 := strings.Index(completeurlpath, "?")

	if index02 == -1 {
		return completeurlpath

	}

	return completeurlpath[index01:index02]

}
func getrequestcompleteurlpath(requestbody string) string {
	index01 := strings.Index(requestbody, " /")
	index02 := strings.Index(requestbody, " HTTP")

	return requestbody[index01+1 : index02]
}

func getquery(completeurlpath string) map[string]string {
	m := make(map[string]string)
	//method := getrequestmethod()

	querystruct := strings.Index(completeurlpath, "?")

	//method == "GET" &&
	if querystruct != -1 {
		rawquery := completeurlpath[querystruct+1:]
		if rawquery != "" {
			pars := strings.Split(rawquery, "&")
			for _, par := range pars {
				parkv := strings.Split(par, "=")
				m[parkv[0]] = parkv[1] // 等号前面是key,后面是value
			}
			return m
		}
		return nil

		//	method == "GET" &&
	} else if querystruct == -1 {
		return nil
	}
	return nil
}

func postquery(httpbody string) map[string]string {
	m := make(map[string]string)
	index01 := strings.Index(httpbody, "\n\n")
	postparameter := httpbody[index01+2:]

	//fmt.Println(postparameter)

	if postparameter != "" {
		pars := strings.Split(postparameter, "&")
		for _, par := range pars {
			parkv := strings.Split(par, "=")
			m[parkv[0]] = parkv[1] // 等号前面是key,后面是value
		}
		return m
	}
	return nil

}

func (c *HTTPStruct) GET(urlpath string, handler Handler) {
	Handlers[urlpath] = handler

}

func (c *HTTPStruct) POST(urlpath string, handler Handler) {
	Handlers[urlpath] = handler
}

func (c *HTTPStruct) Middleware(next func(context *Context)) Handler {
	return func(context *Context) {
		fmt.Println("函数调用之前")
		next(context)
		fmt.Println("函数调用之后")
	}
}
