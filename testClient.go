/*
@Time : 19/4/25 下午8:30 
@Author : gongjiapeng
@File : testClient
@Software: GoLand
*/
package main

import (
	"time"
	"net"
	"fmt"
	"goconnpool/poolConn"
)
//测试tcp链接
func main() {
	pool,_:= poolConn.InitConnPool(10,time.Second*10, func() (poolConn.ConnRes, error) {
		return net.Dial("tcp",":8000")
	})
	buf := make([]byte,1024)
	//拿连接池长度
	l := pool.Len()
	fmt.Println(l)
	//获取一个连接
	conn1 ,_:=pool.Get()
	conn1.(net.Conn).Write([]byte("hello"))
	n,err:=conn1.(net.Conn).Read(buf)
	if err!=nil {
		return
	}
	fmt.Println("返回值：",string(buf[0:n]))

	conn2 ,_:=pool.Get()
	conn2.(net.Conn).Write([]byte("world"))
	n,errs:=conn2.(net.Conn).Read(buf)
	if errs!=nil {
		return
	}
	fmt.Println("返回值：",string(buf[0:n]))
}