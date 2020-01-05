package process

import (
	"encoding/json"
	"fmt"
	"net"
	"sggStudyGo/chatroom/common/message"
	"sggStudyGo/chatroom/server/model"
	"sggStudyGo/chatroom/server/utils"
)

type UserProcess struct {
	Conn net.Conn
	UserId int // 表示当前连接所属哪个用户
}

// 处理登陆操作
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 先从代码中取出mes中的data数据，并json反序列化
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err: ", err)
		return
	}

	// 先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 声明一个 loginResMes
	var loginResMes message.LoginResMes

	// 去redis中验证用户是否存在
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 300
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}

	} else {
		loginResMes.Code = 200
		// 登陆成功，需要记录在线用户
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		// 将当前的userid放入到loginResMes.UserId中
		for id := range userMgr.GetAllOnlineUser() {
			loginResMes.UserIds = append(loginResMes.UserIds, id)
		}

		// 通知其他在线用户该用户已上线
		this.NotifyOthersOnlineUser(loginMes.UserId)

		fmt.Println("登陆成功", user)
	}


	// 将loginResMes转为json
	loginResMesByte, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) err: ", err)
		return
	}

	// 将loginResMesByte转换为string类型赋值给resMes.Data
	resMes.Data = string(loginResMesByte)

	// 将 resMes 转为json
	resMesByte, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) err: ", err)
		return
	}

	// 发送data 封装了writePkg用来发送数据
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(resMesByte)

	return
}


// 处理注册操作
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 先从代码中取出mes中的data数据，并json反序列化
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err: ", err)
		return
	}

	// 先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	// 声明一个 loginResMes
	var registerResMes message.RegisterResMes

	// 去redis中验证用户是否存在
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册出现的未知错误"
		}

	} else {
		registerResMes.Code = 200
	}

	// 将registerResMes转为json
	registerResMesByte, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal(registerResMes) err: ", err)
		return
	}

	// 将loginResMesByte转换为string类型赋值给resMes.Data
	resMes.Data = string(registerResMesByte)

	// 将 resMes 转为json
	resMesByte, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) err: ", err)
		return
	}

	// 发送data 封装了writePkg用来发送数据
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(resMesByte)

	return
}

func (this *UserProcess) NotifyOthersOnlineUser(userid int) {
	// 遍历 onlineUsers，然后一个一个的发送，NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		if id == userid {
			continue
		}

		// 开始通知
		up.NotifyMeOnline(userid)
	}
}

func (this *UserProcess) NotifyMeOnline(userid int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userid
	notifyUserStatusMes.Status = message.UserOnline

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal(notifyUserStatusMes) err: ", err)
		return
	}

	mes.Data = string(data)

	// 将mes序列化
	mesByte, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) err: ", err)
		return
	}

	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	err = tf.WritePkg(mesByte)
	if err != nil {
		fmt.Println("tf.WritePkg(mesByte) err: ", err)
		return
	}
}
