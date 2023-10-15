package main

import (
	"fmt"
	"github.com/idun886/MYHTTP"
	"net/http"
)

func main() {
	myhttp := MYHTTP.New()
	//myhttp := MYHTTP.HTTPStruct{}

	myhttp.GET("/hello", myhttp.Middleware(func(context *MYHTTP.Context) {
		context.JSON(http.StatusOK, MYHTTP.H{
			"200": "hello",
		})
	}))

	//myhttp.GET("/hello", func(context *MYHTTP.Context) {
	//	//content := fmt.Sprintf("HTTP/1.1 200\n%s\n\n%s", "Content-Type: text/plain; charset=utf-8", "你好")
	//	//context.Conn.Write([]byte(content))
	//	newmap := map[string]string{"test": "test01"}
	//	context.AddHeadler("test", "test01")
	//	context.AddHeadler("test03", "test04")
	//	//context.WriteString(http.StatusOK, "你好")
	//
	//	query, err := context.GetQuery("test")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(query)
	//
	//	context.JSON(http.StatusOK, MYHTTP.H{
	//		"code": "200",
	//		"map":  newmap,
	//	})
	//})

	myhttp.POST("/post", func(context *MYHTTP.Context) {

		context.JSON(http.StatusOK, MYHTTP.H{
			"code": "200",
			"map":  "map1",
		})
		fmt.Println(context.PostQuery["test"])
		//content := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n%s", "post hello")
		//context.Conn.Write([]byte(content))
	})

	myhttp.MainConnetc("8000")
}
