package process

import (
	"encoding/json"
	"fmt"
	"sggStudyGo/chatroom/client/utils"
	"sggStudyGo/chatroom/common/message"
)

type SmsProcess struct {

}

// 发送群聊的消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	// 创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	// 创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content // 内容
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	// 序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail = ", err.Error())
		return
	}

	mes.Data = string(data)

	// 再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail = ", err.Error())
		return
	}

	// 将mes发送
	tf := utils.Transfer{
		Conn:CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes tf.WritePkg(data) err = ", err.Error())
		return
	}

	return
}
