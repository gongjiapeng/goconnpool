/*
@Time : 19/4/25 下午5:52 
@Author : gongjiapeng
@File : server.go
@Software: GoLand
*/
package poolConn

import (
	"time"
	"sync"
	"errors"
)

//接口
type ConnRes interface {
	Close() error
}
//工厂
type Factory func()(ConnRes,error)
//链接
type Conn struct {
	conn ConnRes
	time time.Time
}
//链接池
type ConnPoll struct {
	conn chan *Conn
	mu sync.Mutex//锁
	fac Factory
	closed bool
	connTimeOut time.Duration
}

func InitConnPool(poolCount int,timeOut time.Duration,fac Factory)(*ConnPoll,error) {
	if poolCount <0 {
		return nil ,errors.New("链接数不能为空")
	}
	if timeOut <0{
		return nil,errors.New("超时时间不能为空")
	}
	pool := &ConnPoll{
		mu:sync.Mutex{},
		closed:false,
		connTimeOut:timeOut,
		fac:fac,
		conn:make(chan *Conn,poolCount),
	}
	for i:=0; i<poolCount;i++  {
		connRes ,err:= pool.fac()
		if err != nil {
			//关闭
			pool.Close()
		}
		pool.conn <- &Conn{conn:connRes,time:time.Now()}
	}
	return pool,nil
}
//获取链接
func (pool *ConnPoll) Get()( ConnRes, error) {
	if pool.closed {
		return nil,errors.New("链接池以关闭")
	}
	for  {
		select {
		case conn,ok := <-pool.conn:
			if !ok {
				return nil,errors.New("链接池以关闭")
			}
			if time.Now().Sub(conn.time)>pool.connTimeOut {
				conn.conn.Close()
				continue
			}

			return conn.conn ,nil
		default:
			//如果没有重新创建
			connRes ,err := pool.fac()
			if err!=nil {
				return nil,err
			}
			return connRes,nil
		}
	}

}
//放回连接池
func (pool *ConnPoll) Put(conn ConnRes) error {
	if pool.closed {
		return errors.New("连接池已关闭")
	}
	select {
	case pool.conn<-&Conn{conn:conn,time:time.Now()}:
		return nil
	default:
		conn.Close()
		return errors.New("连接池以满")
	}

}
//关闭
func (pool *ConnPoll) Close()  {
	if pool.closed {
		return
	}
	//线程安全 加锁
	pool.mu.Lock()
	pool.closed=true
	//通道关闭
	close(pool.conn)
	//通道里面的链接关闭
	for conn:=range pool.conn{
		conn.conn.Close()
	}
	pool.mu.Unlock()
}
//返回连接池的长度
func (poll *ConnPoll) Len() int {
	return len(poll.conn)
}


