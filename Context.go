package MYHTTP

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Context struct {
	Conn           net.Conn //传入链接
	RequestMethod  string   //请求方法
	RequestHost    string   //请求域名
	RequestUrlPath string   //请求url路径

	getQuery  map[string]string //get请求参数
	PostQuery map[string]string //post请求参数

	WriteBody string //写入时的数据
	//Headler      string //头部
	//ResponseBody string
}

func (c *Context) GetQuery(key string) (string, error) {

	err := errors.New("Context class 's GetQuery function error")
	query, ok := c.getQuery[key]
	if ok {
		return query, nil
	}
	return "", err
}

//func (c *Context) PostQuery(key string) (string, error) {
//	err := errors.New("Context class 's PostQuery function error")
//
//}

// 响应实例
// HTTP/1.1 200 OK
//
// hello
//var response_row string = "HTTP/1.1 200\nContent-Type: text/plain; charset=utf-8"

func (c *Context) AddHeadler(key, value string) string {
	//var response_row string = "HTTP/1.1 200\nContent-Type: text/plain; charset=utf-8"
	headler := fmt.Sprintf("%s: %s", key, value)

	c.WriteBody = fmt.Sprintf(c.WriteBody+"\n%s", headler)

	//fmt.Println(c.WriteBody)
	return c.WriteBody

}

func (c *Context) WriteString(code int, body string) {
	//var response_row string = c.AddHeadler("", "")
	codestr := strconv.Itoa(code)
	c.WriteBody = strings.Replace(c.WriteBody, "200", codestr, 1)
	c.WriteBody = fmt.Sprintf(c.WriteBody+"\n\n%s", body)

	c.Conn.Write([]byte(c.WriteBody))

}

func (c *Context) JSON(code int, jsonmap any) {
	jsonstr, err := json.Marshal(jsonmap)
	if err != nil {
		fmt.Printf("Context class 's JSON function error:%s", err)
	}

	codestr := strconv.Itoa(code)
	c.WriteBody = strings.Replace(c.WriteBody, "200", codestr, 1)
	c.WriteBody = fmt.Sprintf(c.WriteBody+"\n\n%s", jsonstr)
	c.Conn.Write([]byte(c.WriteBody))
}
