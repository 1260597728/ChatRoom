package main

import (
	"fmt"
	"net"
	"sggStudyGo/chatroom/server/model"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close()

	// 创建总控并调用
	pro := &Processor{
		Conn: conn,
	}
	err := pro.MainProcess()
	if err != nil {
		fmt.Println("pro.MainProcess err: ", err)
		return
	}
}

func init() {
	// 初始化redis连接池
	// 先调用initPool生成redis连接池
	initPool("127.0.0.1:6379", 16, 0, 100 * time.Second)

	// initPool初始化连接池之后initUserDao会用到
	initUserDao()
}

// 完成对UserDao的初始化任务
func initUserDao() {
	// 这里的pool 本身就是一个全局的变量
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {

	fmt.Println("服务器8889端口正在监听……")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Lister err: ", err)
		return
	}

	// 一旦监听成功，就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端来连接服务器……")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err: ", err)
			return
		}

		// 连接成功，就启动一个协程与客户端保持连接通讯
		go process(conn)
	}

}
