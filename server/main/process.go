package main

import (
	"fmt"
	"io"
	"net"
	"sggStudyGo/chatroom/common/message"
	process2 "sggStudyGo/chatroom/server/process"
	"sggStudyGo/chatroom/server/utils"
)

type Processor struct {
	Conn net.Conn
}

// 根据客户端发送的类型不同，决定调用哪个函数处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType: // 登陆操作
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType: // 注册操作
		up := &process2.UserProcess{
			Conn:this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在，无法处理")
	}

	return
}

// 总控
func (this *Processor) MainProcess() (err error) {
	// 循环获取客户端发送的消息
	for {
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出")
			} else {
				fmt.Println("tf.ReadPkg err: ", err)
			}
			return err
		}

		// 根据类型不同处理不同操作
		err = this.serverProcessMes(&mes)
		if err != nil {
			fmt.Println("serverProcessMes err: ", err)
			return err
		}

		fmt.Println(mes)
	}
}
