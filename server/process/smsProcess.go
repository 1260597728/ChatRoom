package process

import (
	"encoding/json"
	"fmt"
	"net"
	"sggStudyGo/chatroom/common/message"
	"sggStudyGo/chatroom/server/utils"
)

type SmsProcess struct {

}

func (this *SmsProcess) SendGroupMes(mes *message.Message) {

	// 取出mes的内容
	var smsRes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsRes)
	if err != nil {
		fmt.Println("json.Unmarshal err: ", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) err: ", err)
		return
	}


	// 遍历服务器端的onlineUsers map[int]*UserProcess
	// 将消息转发出去
	for id, up := range userMgr.onlineUsers {
		// 不需要发给自己
		if id == smsRes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}
}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	tf := utils.Transfer{
		Conn:conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("tf.WritePkg(data) err: ", err)
		return
	}
}
