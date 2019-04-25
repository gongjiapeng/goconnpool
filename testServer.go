 package main

 import (
	 "net"
	 "fmt"
	 "io"
 )

 func main() {
	 lis,err:=net.Listen("tcp",":8000")
	 if err!=nil {
		 return
	 }
	 for  {
		 conn,err:=lis.Accept()
		 if err!=nil {
			 fmt.Println(err)
			 continue
		 }
		 go handleConn(conn)
	 }
 }
 func handleConn(c net.Conn)  {
	 defer c.Close()
	 buf:=make([]byte,1024)
	 for  {
	 	 n,err:=c.Read(buf)
		 if err!=nil {
			 return
		 }
		 _, errs := io.WriteString(c,"拿到了")
		 if errs != nil {
			 return // e.g., client disconnected
		 }
		 fmt.Println("conn read",string(buf[0:n]))
	 }
 }