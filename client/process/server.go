package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sggStudyGo/chatroom/client/utils"
	"sggStudyGo/chatroom/common/message"
)

// 登陆成功，显示在线页面
func ShowMenu(username string) {
	fmt.Println("------你好，" + username + "-------")
	fmt.Println("------1. 显示在线用户列表-------")
	fmt.Println("------2. 发送信息-------")
	fmt.Println("------3. 信息列表-------")
	fmt.Println("------4. 退出系统-------")
	fmt.Println("请选择：")

	// 获取用户终端输入
	var key int

	// 需要发送的内容
	var content string
	smsProcess := &SmsProcess{}

	for true {
		// 这个得放在for循环里面，放在外面获取就成死循环了，就一直是这个值
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			//fmt.Println("显示在线用户列表")
			outputOnlineUser()
		case 2:
			fmt.Println("你想对大家说点什么：）")
			fmt.Scanf("%s\n", &content)

			smsProcess.SendGroupMes(content)
		case 3:
			fmt.Println("信息列表")
		case 4:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入的选项不正确……请重新输入")
		}
	}

}

func serverProcessMes(conn net.Conn) {
	// 不停的读服务器端发过来的数据
	tf := &utils.Transfer{
		Conn: conn,
	}
	// conn.Read 没信息过来的话会一直阻塞
	for {
		// 阻塞监听服务器端有没有发送过来数据
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("client tf.ReadPkg err: ", err)
			return
		}

		// 读取到信息
		switch mes.Type {
		case message.NotifyUserStatusMesType: // 有人上线
			// 取出NotifyUserStatus
			// 把这个用户信息，状态保存到客户端中的map[int]User
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			// 把这个用户的信息保存到客户端map[int]User中
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType: // 有人发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")
		}
	}
}
