package process

import (
	"fmt"
	"sggStudyGo/chatroom/client/model"
	"sggStudyGo/chatroom/common/message"
)

// 客户端需要维护的map
var onlineUsers = make(map[int]*message.User, 10)
var CurUser model.CurUser // 在用户登陆成功之后，完成对CurUser的初始化

// 编写一个方法，处理返回的 NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId:     notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status

	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}

// 在客户端显示当前在线用户
func outputOnlineUser() {
	// 遍历一下 onlineUsers
	fmt.Println("当前在线用户列表：")
	for id := range onlineUsers {
		fmt.Println("用户id：", id)
	}
}
