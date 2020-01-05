package main

import (
	"fmt"
	"os"
	"sggStudyGo/chatroom/client/process"
)

var (
	userId  int    // 用户id
	userPwd string // 用户密码
	userName string // 用户昵称
)

func main() {

	// 终端输入记录用户选择
	var choose int

	for true {
		fmt.Println("------ 欢迎登陆用户系统 ------")
		fmt.Printf("\t\t\t\n 1 登陆聊天室")
		fmt.Printf("\t\t\t\n 2 注册用户")
		fmt.Printf("\t\t\t\n 3 退出系统")
		fmt.Printf("\n请选择：")

		// 终端获取用户输入
		fmt.Scanln(&choose)
		switch choose {
		case 1:
			fmt.Println("### 登陆聊天室 ###")

			// 获取终端用户输入id和密码
			getUserScan()

			// 执行登陆操作
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("### 注册用户 ###")

			// 获取终端用户输入id和密码
			getUserScan()
			fmt.Println("请输入昵称：")
			fmt.Scanf("%s\n", &userName)

			// 执行注册操作
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
			// 不管成功还是失败，直接退出
			os.Exit(0)
		case 3:
			fmt.Println("### 退出系统 ###")
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")
		}
	}
}

// 获取用户终端输入id和密码
func getUserScan() {
	fmt.Println("请输入id：")
	//fmt.Scanln(&userId)
	fmt.Scanf("%d\n", &userId)
	fmt.Println("请输入密码：")
	//fmt.Scanln(&userPwd)
	fmt.Scanf("%s\n", &userPwd)
}
