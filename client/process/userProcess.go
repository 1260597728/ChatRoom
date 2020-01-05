package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sggStudyGo/chatroom/client/utils"
	"sggStudyGo/chatroom/common/message"
	"strconv"
)

type UserProcess struct {
}

// 用户注册
func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	// 连接到服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return
	}

	// 定义消息结构体
	var messageStruct message.Message
	// 消息类型
	messageStruct.Type = message.RegisterMesType

	// 定义消息注册结构体
	var registerStruct message.RegisterMes
	registerStruct.User.UserId = userId
	registerStruct.User.UserPwd = userPwd
	registerStruct.User.UserName = userName

	// 将loginStruct序列化
	loginStructByte, err := json.Marshal(registerStruct)
	if err != nil {
		fmt.Println("json.Marshal registerStruct err: ", err)
		return
	}

	// 消息内容
	messageStruct.Data = string(loginStructByte)

	// 将消息结构体进行序列化
	messageStructByte, err := json.Marshal(messageStruct)
	if err != nil {
		fmt.Println("json.Marshal messageStruct err: ", err)
		return
	}

	// 处理服务器返回的信息
	tf := &utils.Transfer{
		Conn: conn,
	}
	// 发送信息给服务器端
	err = tf.WritePkg(messageStructByte)
	if err != nil {
		fmt.Println("tf.WritePkg(messageStructByte) err: ", err)
		return
	}

	// 从服务器端接收消息
	// 这里的mes 指的是 RegisterResMes
	mes, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err: ", err)
		return
	}

	// 将mes.data 的部分反序列化成 RegisterResMes
	var RegisterMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &RegisterMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &RegisterMes) err: ", err)
		return
	}

	if RegisterMes.Code == 200 {
		fmt.Println("注册成功")
	} else {
		fmt.Println(RegisterMes.Error)
	}

	return
}

// 关联一个用户登陆的方法
func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	// 连接到服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return
	}

	// 定义消息结构体
	var messageStruct message.Message
	// 消息类型
	messageStruct.Type = message.LoginMesType

	// 定义消息登陆结构体
	var loginStruct message.LoginMes
	loginStruct.UserId = userId
	loginStruct.UserPwd = userPwd

	// 将loginStruct序列化
	loginStructByte, err := json.Marshal(loginStruct)
	if err != nil {
		fmt.Println("json.Marshal loginStruct err: ", err)
		return
	}

	// 消息内容
	messageStruct.Data = string(loginStructByte)

	// 将消息结构体进行序列化
	messageStructByte, err := json.Marshal(messageStruct)
	if err != nil {
		fmt.Println("json.Marshal messageStruct err: ", err)
		return
	}

	// 将数据发送给服务器
	// 先将数据大小发送给服务器，将数据大小转为[]byte
	pkgLen := uint32(len(messageStructByte))
	// 一个uint32等于4个字节，一字节等于8bit
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) err: ", err)
		return
	}
	//fmt.Printf("客户端发送消息的长度: %d, 内容是：%s", len(messageStructByte), string(messageStructByte))

	// 发送消息本身
	_, err = conn.Write(messageStructByte)
	if err != nil {
		fmt.Println("conn.Write(messageStructByte) err: ", err)
		return
	}

	// 休眠10s
	//time.Sleep(10 * time.Second)

	// 处理服务器返回的信息
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err: ", err)
		return
	}

	// 将mes.data 的部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &loginResMes) err: ", err)
		return
	}

	if loginResMes.Code == 200 {
		//fmt.Println("登陆成功")

		// 初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UserIds {

			// 可以判断下，当前登陆用户可以不用展示
			if v == userId {
				continue
			}

			fmt.Println("用户id：", v)

			// 完成客户端的 onlineUsers 完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println()

		// 登陆成功之后还需要在客户端启动一个协程
		// 该协程保持和服务器端的通讯
		// 如果服务器有数据推送给客户端，则接收并显示在客户端的终端
		go serverProcessMes(conn)

		// 循环展示登陆成功之后的菜单
		//username := fmt.Sprintf("%s", userId)
		username := strconv.Itoa(userId)
		for {
			ShowMenu(username)
		}

	} else {
		fmt.Println(loginResMes.Error)
	}

	return
}
